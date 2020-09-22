package apple

import (
	"context"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"strings"
)

const redisKeyPrefix = "iap:receipt:"

func getKafkaReader(addr []string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  addr,
		GroupID:  "iap-verifier",
		Topic:    Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

type Consumer struct {
	vrfClient  VerificationClient
	subsClient SubsClient
	rdb        *redis.Client
	reader     *kafka.Reader
	logger     *zap.Logger
	ctx        context.Context
}

func NewConsumer(prod bool, logger *zap.Logger) Consumer {
	return Consumer{
		vrfClient:  NewVerificationClient(prod),
		subsClient: NewSubsClient(prod),
		rdb: redis.NewClient(&redis.Options{
			Addr:     config.MustRedisAddress().Pick(prod),
			Password: "",
			DB:       0,
		}),
		reader: getKafkaReader(config.MustKafkaAddress().PickSlice(prod)),
		logger: logger,
		ctx:    context.Background(),
	}
}

func (c Consumer) Consume() {
	defer c.logger.Sync()
	sugar := c.logger.Sugar()

	sugar.Info("Start consuming...")

	for {
		m, err := c.reader.ReadMessage(c.ctx)
		if err != nil {
			sugar.Error(err)
			continue
		}

		sugar.Infof("message at topic: %v, partition %v, offset %v, %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		// TODO:
		// 1. get receipt from redis; then fallback to subscription api
		// 2. Verify receipt
		// 3. forward response to subscription api
	}
}

// Turn string like `180000520108935,Production` into `iap:receipt:180000520108935-Production`
func receiptRedisKey(name string) string {
	return redisKeyPrefix + strings.Join(strings.Split(name, ","), "-")
}

func getOrigTxID(name string) string {
	return strings.Split(name, ",")[0]
}

func (c Consumer) getReceipt(name string) (string, error) {
	val, err := c.rdb.Get(c.ctx, receiptRedisKey(name)).Result()
	if err == nil {
		return val, nil
	}

	receipt, err := c.subsClient.GetReceipt(getOrigTxID(name))
	if err != nil {
		return "", err
	}

	return receipt, nil
}

func (c Consumer) verify(name string) error {
	defer c.logger.Sync()
	sugar := c.logger.Sugar()

	receipt, err := c.getReceipt(name)
	if err != nil {
		sugar.Error(err)
		return err
	}

	body, err := c.vrfClient.Verify(receipt)
	if err != nil {
		sugar.Error(err)
		return err
	}

	resp, errs := c.subsClient.SaveReceipt(body)
	if errs != nil {
		sugar.Error(errs)
		return errs[0]
	}

	sugar.Infof("Saved receipt. Resp %d", resp.StatusCode)

	return nil
}
