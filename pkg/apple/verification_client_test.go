package apple

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"testing"
)

func TestVerificationClient_Verify(t *testing.T) {
	config.MustSetupViper()

	subsAPI := NewSubsClient(false)

	r, err := subsAPI.GetReceipt("1000000619244062")
	if err != nil {
		t.Error(err)
	}

	vrfClient := NewVerificationClient(false)
	resp, body, errs := vrfClient.Verify(r)
	if errs != nil {
		t.Error(errs)
	}

	t.Logf("Response status %d", resp.StatusCode)
	t.Logf("%s", body)
}
