package apple

import (
	"encoding/json"
	"errors"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/tidwall/gjson"
	"log"
)

type SubsClient struct {
	key     string
	baseURL string
}

func NewSubsClient(prod bool) SubsClient {

	return SubsClient{
		key:     config.MustAPIKey().Pick(prod),
		baseURL: config.MustAPIBaseURL().Pick(prod),
	}
}

// GetReceipt tries to get a receipt file from various API.
// This is used as a fallback when the receipt cannot be found in redis.
func (c SubsClient) GetReceipt(origTxID string) (string, error) {
	url := c.baseURL + "/apple/receipt/" + origTxID

	resp, b, errs := fetch.New().Get(url).SetBearerAuth(c.key).EndBytes()
	if errs != nil {
		return "", errs[0]
	}

	if resp.StatusCode >= 400 {
		var respErr render.ResponseError
		if err := json.Unmarshal(b, &respErr); err != nil {
			return "", err
		}

		return "", &respErr
	}

	result := gjson.GetBytes(b, "receipt")

	if !result.Exists() {
		return "", errors.New("receipt not found from subscription api")
	}

	return result.String(), nil
}

// VerifyReceipt send a receipt to subscription api to get
// Subscription response.
// Treat http status code above 400 as error.
func (c SubsClient) VerifyReceipt(receipt string) ([]byte, error) {
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
