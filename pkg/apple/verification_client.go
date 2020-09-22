package apple

import (
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/fetch"
	"io"
	"log"
)

type VerificationClient struct {
	production bool
	sandboxUrl string
	prodUrl    string
	password   string
}

func NewVerificationClient(prod bool) VerificationClient {
	return VerificationClient{
		production: prod,
		sandboxUrl: "https://sandbox.itunes.apple.com/verifyReceipt",
		prodUrl:    "https://buy.itunes.apple.com/verifyReceipt",
		password:   config.MustIAPSecret(),
	}
}

func (c VerificationClient) pickUrl() string {
	if c.production {
		log.Print("Using IAP production url")
		return c.prodUrl
	}

	log.Print("Using IAP sandbox url")
	return c.sandboxUrl
}

func (c VerificationClient) Verify(receipt string) (io.ReadCloser, error) {
	payload := VerificationPayload{
		ReceiptData:            receipt,
		Password:               c.password,
		ExcludeOldTransactions: false,
	}

	resp, errs := fetch.New().
		Post(c.pickUrl()).
		SendJSON(payload).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp.Body, nil
}
