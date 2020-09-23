package apple

import (
	"context"
	"fmt"
	"testing"

	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/robfig/cron/v3"
	"github.com/segmentio/kafka-go"
)

func TestNewKafkaWriter(t *testing.T) {
	config.MustSetupViper()

	writer := NewKafkaWriter(config.MustKafkaAddress().PickSlice(false))

	for i := 0; i < 10; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", i)),
			Value: []byte(fmt.Sprintf("Hello world %d", i)),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestCron(t *testing.T) {
	c := cron.New()

	_, err := c.AddFunc("1-59 * * * *", func() {
		println("Hello")
	})

	if err != nil {
		t.Error(err)
	}

	c.Start()

	for {
	}
}
