package apple

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"go.uber.org/zap/zaptest"
	"path/filepath"
	"testing"
)

func TestWorker_Verify(t *testing.T) {
	config.MustSetupViper()

	w := NewWorker(false, zaptest.NewLogger(t))

	filename := filepath.Join(mustHomeDir(), "receipt/user-id/5/5a0a1a22-505f-4c93-bcda-d52d1db3fa8a.log")

	s, err := w.Verify(filename)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", s)
}
