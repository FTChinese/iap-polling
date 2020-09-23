package apple

import (
	"context"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/db"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
)

const (
	Topic     = "iap-polling"
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

type Producer struct {
	db         *sqlx.DB
	vrfClient  VerificationClient
	subsClient SubsClient
	rdb        *redis.Client
	writer     *kafka.Writer
	logger     *zap.Logger
	ctx        context.Context
}

func NewProducer(prod bool, logger *zap.Logger) *Producer {
	return &Producer{
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

func (p *Producer) Produce() {
	defer p.logger.Sync()
	sugar := p.logger.Sugar()

	rows, err := p.db.Queryx(StmtSubs)
	if err != nil {
		return
	}

	subs := Subscription{}
	for rows.Next() {
		err := rows.StructScan(&subs)
		if err != nil {
			sugar.Error(err)
			continue
		}

		sugar.Infof("%#v\n", subs)

		err = p.writer.WriteMessages(
			context.Background(),
			kafka.Message{
				Key:   []byte(subs.OriginalTransactionID),
				Value: []byte(subs.String()),
			},
		)

		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}
}

func (p *Producer) getReceipt(name string) (string, error) {
	val, err := p.rdb.Get(p.ctx, receiptRedisKey(name)).Result()
	if err == nil {
		return val, nil
	}

	receipt, err := p.subsClient.GetReceipt(getOrigTxID(name))
	if err != nil {
		return "", err
	}

	return receipt, nil
}

func (p *Producer) verify(name string) error {
	defer p.logger.Sync()
	sugar := p.logger.Sugar()

	receipt, err := p.getReceipt(name)
	if err != nil {
		sugar.Error(err)
		return err
	}

	body, err := p.vrfClient.Verify(receipt)
	if err != nil {
		sugar.Error(err)
		return err
	}

	resp, errs := p.subsClient.SaveReceipt(body)
	if errs != nil {
		sugar.Error(errs)
		return errs[0]
	}

	sugar.Infof("Saved receipt. Resp %d", resp.StatusCode)

	return nil
}

func (p *Producer) Close() {
	log.Print("Closing db and kafka producer")
	if err := p.writer.Close(); err != nil {
		log.Fatal(err)
	}

	if err := p.db.Close(); err != nil {
		log.Fatal(err)
	}
}
