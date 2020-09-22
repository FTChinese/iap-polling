package apple

import (
	"context"
	"github.com/FTChinese.com/iap-polling/pkg/db"
	"github.com/FTChinese/go-rest/connect"
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
	db     *sqlx.DB
	writer *kafka.Writer
	logger *zap.Logger
}

func NewProducer(conn connect.Connect, addr []string, logger *zap.Logger) *Producer {
	return &Producer{
		db:     db.MustNewDB(conn),
		writer: NewKafkaWriter(addr),
		logger: logger,
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

func (p *Producer) Close() {
	log.Print("Closing db and kafka producer")
	if err := p.writer.Close(); err != nil {
		log.Fatal(err)
	}

	if err := p.db.Close(); err != nil {
		log.Fatal(err)
	}
}
