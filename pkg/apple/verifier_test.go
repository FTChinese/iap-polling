package apple

import (
	"context"
	"fmt"
	"go.uber.org/zap/zaptest"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/segmentio/kafka-go"
)

func mustHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return h
}

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

func TestVerifier_getReceipt(t *testing.T) {
	v := NewVerifier(false, zaptest.NewLogger(t))

	r, err := v.getReceipt(Subscription{
		Environment:           EnvSandbox,
		OriginalTransactionID: "1000000619244062",
	})

	if err != nil {
		t.Error(err)
	}

	t.Logf("Receipt: %s", r)
}

func TestVerifier_Produce(t *testing.T) {
	v := NewVerifier(false, zaptest.NewLogger(t))

	b, err := ioutil.ReadFile(filepath.Join(mustHomeDir(), "config/apple_verified_receipt.json"))
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s", b)

	err = v.Produce("test", b)
	if err != nil {
		t.Error(err)
	}
}
