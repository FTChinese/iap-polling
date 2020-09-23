package apple

import (
	"context"
	"errors"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/db"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Topic     = "iap-polled-receipt"
	Partition = 0
)

func NewKafkaWriter(addr []string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  addr,
		Topic:    Topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
}

type Verifier struct {
	db         *sqlx.DB
	vrfClient  VerificationClient
	subsClient SubsClient
	rdb        *redis.Client
	writer     *kafka.Writer
	logger     *zap.Logger
	ctx        context.Context
}

func NewVerifier(prod bool, logger *zap.Logger) *Verifier {
	return &Verifier{
		db:         db.MustNewDB(config.MustDBConn(prod)),
		vrfClient:  NewVerificationClient(prod),
		subsClient: NewSubsClient(prod),
		rdb: redis.NewClient(&redis.Options{
			Addr:     config.MustRedisAddress().Pick(prod),
			Password: "",
			DB:       0,
		}),
		writer: NewKafkaWriter(config.MustKafkaAddress().PickSlice(prod)),
		logger: logger,
		ctx:    context.Background(),
	}
}

func (v *Verifier) LoadSubs(ch chan<- Subscription) error {
	defer v.logger.Sync()
	sugar := v.logger.Sugar()

	defer close(ch)

	rows, err := v.db.Queryx(StmtSubs)
	if err != nil {
		sugar.Error(err)
		return err
	}

	subs := Subscription{}
	for rows.Next() {
		err := rows.StructScan(&subs)
		if err != nil {
			sugar.Error(err)
			continue
		}

		sugar.Infof("%#v\n", subs)

		ch <- subs
	}

	return nil
}

// getReceipt tries to load a receipt for a Subscription from redis, and fallback to
// subscription api if not found.
func (v *Verifier) getReceipt(subs Subscription) (string, error) {
	val, err := v.rdb.Get(v.ctx, subs.ReceiptKeyName()).Result()
	if err == nil {
		return val, nil
	}

	receipt, err := v.subsClient.GetReceipt(subs.OriginalTransactionID)
	if err != nil {
		return "", err
	}

	return receipt, nil
}

// Verify performs verification by a receipt's original transaction id.
func (v *Verifier) Verify(subs Subscription) error {
	defer v.logger.Sync()
	sugar := v.logger.Sugar()

	// Get existing receipt first
	receipt, err := v.getReceipt(subs)
	if err != nil {
		sugar.Error(err)
		return err
	}

	// Verify the receipt against app store.
	resp, errs := v.vrfClient.Verify(receipt)
	if errs != nil {
		sugar.Error(err)
		return errs[0]
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sugar.Error("App store response status code %d", resp.StatusCode)
		return errors.New("app store response not ok")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Send the response data to kafka as is.
	err = v.Produce(subs.OriginalTransactionID, body)

	if err != nil {
		sugar.Errorf("failed to write messages:", err)
	}

	return nil
}

func (v *Verifier) Produce(key string, body []byte) error {
	// Send the response data to kafka as is.
	return v.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: body,
		},
	)
}

func (v *Verifier) Close() {
	log.Print("Closing db and kafka producer")
	if err := v.writer.Close(); err != nil {
		log.Fatal(err)
	}

	if err := v.db.Close(); err != nil {
		log.Fatal(err)
	}
}
