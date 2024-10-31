package shipment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Shipment is the client utility for EVM Shipment
type Shipment struct {
	Client *http.Client
	URL    string
}

//HTTPRequest is a request http client to external domain in 1 host
func (c *Shipment) HTTPRequest(payload []byte) (err error) {
	urlTransaction := "/v2/shipment/tracking/register"
	url := fmt.Sprintf("%s%s", c.URL, urlTransaction)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		errorReponse := ErrorReponse{}
		json.Unmarshal(body, &errorReponse)
		return fmt.Errorf("%s", errorReponse.Error)
	}

	return
}
