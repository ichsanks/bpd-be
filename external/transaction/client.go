package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Transaction is the client utility for Transaction Domain
type Transaction struct {
	Client *http.Client
	URL    string
}

//HTTPRequest is a request http client to external domain in 1 host
func (c *Transaction) HTTPRequest(userID string, IDOpd string, payload []byte) (trxResponse CreateTransactionFormatResponse, err error) {
	urlTransaction := "/v1/transaction/order"
	url := fmt.Sprintf("%s%s", c.URL, urlTransaction)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	httpReq.Header.Set("X-Imers-Tenant-ID", IDOpd)
	httpReq.Header.Set("X-Imers-User-ID", userID)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return trxResponse, err
	}

	if resp.StatusCode != http.StatusOK {
		errorReponse := ErrorReponse{}
		json.Unmarshal(body, &errorReponse)
		return trxResponse, fmt.Errorf("%s", errorReponse.Error)
	}

	json.Unmarshal(body, &trxResponse)

	return
}
