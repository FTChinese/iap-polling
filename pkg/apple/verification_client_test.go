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
	body, err := vrfClient.Verify(r)
	if err != nil {
		t.Error(err)
	}

	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s", b)
}
