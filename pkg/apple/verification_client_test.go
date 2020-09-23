package apple

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"io/ioutil"
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
	resp, errs := vrfClient.Verify(r)
	if errs != nil {
		t.Error(errs)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s", b)
}
