package apple

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestLinkMigration_saveLinkErrLog(t *testing.T) {
	config.MustSetupViper()

	m := NewLinkMigration(false, zaptest.NewLogger(t))

	err := m.saveLinkErrLog(LinkErrLog{
		LinkInput: LinkInput{
			OriginalTxID: "490000444035809",
			FtcID:        "d815e143-053c-4460-bf2b-91a4854e59a3",
		},
		Field:   "ftcId",
		Code:    "has_valid_non_iap",
		Message: "FTC account already has a valid membership via non-Apple channel",
	})

	if err != nil {
		t.Error(err)
	}
}

func TestLinkMigration_link(t *testing.T) {
	config.MustSetupViper()

	m := NewLinkMigration(false, zaptest.NewLogger(t))

	err := m.link(LinkInput{
		OriginalTxID: "1000000645378944",
		FtcID:        "095c765c-4714-48de-ad61-9be9a72d8e6d",
	})

	if err != nil {
		t.Error(err)
	}
}
