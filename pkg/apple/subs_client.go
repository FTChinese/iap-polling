package apple

import (
	"encoding/json"
	"errors"
	"github.com/FTChinese.com/iap-polling/pkg/config"
	"github.com/FTChinese.com/iap-polling/pkg/fetch"
	"github.com/FTChinese/go-rest/render"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
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

func (c SubsClient) SaveReceipt(body io.ReadCloser) (*http.Response, []error) {
	url := c.baseURL + "/apple/receipt"

	defer body.Close()

	return fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		Send(body).
		End()
}
