package main

import (
	"context"
	"fmt"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func consume() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", config.Topic, config.Partition)

	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)

	b := make([]byte, 10e3)

	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}

		fmt.Println(string(b))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection: ", err)
	}
}
