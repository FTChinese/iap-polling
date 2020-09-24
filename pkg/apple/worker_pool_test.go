package apple

import (
	"context"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/db"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestVerifier_Start(t *testing.T) {
	config.MustSetupViper()

	// Uses IAP production.
	v := &Verifier{
		db:         db.MustNewDB(config.MustDBConn(false)),
		vrfClient:  NewVerificationClient(true),
		subsClient: NewSubsClient(false),
		rdb: redis.NewClient(&redis.Options{
			Addr:     config.MustRedisAddress().Pick(false),
			Password: "",
			DB:       0,
		}),
		writer: NewKafkaWriter(config.MustKafkaAddress().PickSlice(false)),
		logger: zaptest.NewLogger(t),
		ctx:    context.Background(),
	}

	err := v.Start()
	if err != nil {
		t.Error(err)
	}
}
