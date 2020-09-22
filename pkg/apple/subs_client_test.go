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

func TestSubsClient_SaveReceipt(t *testing.T) {
	config.MustSetupViper()

	subsAPI := NewSubsClient(false)

	r, err := subsAPI.GetReceipt("1000000619244062")
	if err != nil {
		t.Error(err)
	}

	vrfClient := NewVerificationClient(false)
	body, err := vrfClient.Verify(r)
	if err != nil {
		t.Error(err)
	}

	resp, errs := subsAPI.SaveReceipt(body)
	if errs != nil {
		t.Error(err)
	}

	t.Logf("Status code %d", resp.StatusCode)
}
