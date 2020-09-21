package main

import (
	"context"
	"fmt"
	"github.com/FTChinese.com/iap-polling/pkg/apple"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/db"
	"github.com/FTChinese/go-rest/connect"
	"github.com/segmentio/kafka-go"
	"log"
)

func Start(conn connect.Connect) {
	myDB := db.MustNewDB(conn)

	w := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	rows, err := myDB.Queryx(apple.StmtSubs)
	if err != nil {
		return
	}

	subs := apple.Subscription{}
	for rows.Next() {
		err := rows.StructScan(&subs)
		if err != nil {
			continue
		}

		fmt.Printf("%#v\n", subs)

		redisKey := subs.ReceiptKeyName()

		err = w.WriteMessages(
			context.Background(),
			kafka.Message{
				Key:   []byte(subs.OriginalTransactionID),
				Value: []byte(redisKey),
			},
		)

		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	if err := myDB.Close(); err != nil {
		log.Fatal("failed to close db:", err)
	}
}
