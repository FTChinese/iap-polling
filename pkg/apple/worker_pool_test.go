package apple

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestVerifier_Start(t *testing.T) {
	config.MustSetupViper()

	v := NewVerifier(false, zaptest.NewLogger(t))

	err := v.Start()
	if err != nil {
		t.Error(err)
	}
}
