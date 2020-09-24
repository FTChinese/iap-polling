package apple

import (
	"context"
	"fmt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"go.uber.org/zap/zaptest"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	config.MustSetupViper()

	v := NewVerifier(false, zaptest.NewLogger(t))

	b, err := ioutil.ReadFile(filepath.Join(mustHomeDir(), "config/apple_verified_receipt.json"))
	if err != nil {
		t.Error(err)
	}

	err = v.Produce("test", b)
	if err != nil {
		t.Error(err)
	}
}

func TestVerifier_Verify(t *testing.T) {
	config.MustSetupViper()

	v := NewVerifier(false, zaptest.NewLogger(t))

	s := Subscription{
		Environment:           EnvSandbox,
		OriginalTransactionID: "1000000619244062",
	}

	err := v.Verify(s)
	if err != nil {
		t.Error(err)
	}
}

func TestVerifier_SaveLog(t *testing.T) {
	config.MustSetupViper()

	v := NewVerifier(false, zaptest.NewLogger(t))

	pl := PollerLog{
		Total:     int64(rand.IntRange(0, 200)),
		Succeeded: int64(rand.IntRange(0, 200)),
		Failed:    int64(rand.IntRange(0, 200)),
		StartUTC:  chrono.TimeNow(),
		EndUTC:    chrono.TimeFrom(time.Now().Add(1 * time.Hour)),
	}

	err := v.SaveLog(&pl)
	if err != nil {
		t.Error(err)
	}
}
