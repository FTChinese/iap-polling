package apple

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"testing"
)

func TestSubsClient_GetReceipt(t *testing.T) {
	config.MustSetupViper()

	client := NewSubsClient(false)

	r, err := client.GetReceipt("1000000619244062")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s", r)
}
