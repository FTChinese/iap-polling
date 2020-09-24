package migrate

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"path/filepath"
	"testing"
)

func TestWorker_Verify(t *testing.T) {
	config.MustSetupViper()

	w := NewWorker(false)

	filename := filepath.Join(mustHomeDir(), "receipt/user-id/5/5a0a1a22-505f-4c93-bcda-d52d1db3fa8a.log")

	s, err := w.Verify(filename)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", s)
}

func TestWorker_SaveMapping(t *testing.T) {
	config.MustSetupViper()

	w := NewWorker(false)

	err := w.SaveMapping(IDMapping{
		TxID:        rand.String(10),
		FtcID:       null.StringFrom(rand.String(10)),
		DeviceToken: null.String{},
		UnionID:     null.String{},
		AbsFilePath: "",
	})

	if err != nil {
		t.Error(err)
	}
}
