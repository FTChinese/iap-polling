package apple

import (
	"encoding/json"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"log"
)

type SubsAPI struct {
	key     string
	baseURL string
}

// NewSubsAPI create a new SubsAPI used to access subscription api.
// If prod is true, visits online production server; otherwise uses development server.
func NewSubsAPI(prod bool) SubsAPI {

	return SubsAPI{
		key:     config.MustAPIKey().Pick(prod),
		baseURL: config.MustAPIBaseURL().Pick(prod),
	}
}

// VerifyReceipt send a receipt to subscription api to get
// Subscription response.
// Treat http status code above 400 as error.
func (c SubsAPI) VerifyReceipt(receipt string) ([]byte, error) {
	url := c.baseURL + "/apple/subs"

	resp, b, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		SendJSON(map[string]string{
			"receiptData": receipt,
		}).
		EndBytes()
	if errs != nil {
		log.Printf("VerifyReceipt: error %v", errs)
		return nil, errs[0]
	}

	if resp.StatusCode >= 400 {

		var respErr render.ResponseError
		if err := json.Unmarshal(b, &respErr); err != nil {
			return nil, err
		}

		log.Printf("VerifyReceipt: subscription api reponse error %v", respErr)
		return nil, &respErr
	}

	return b, nil
}
