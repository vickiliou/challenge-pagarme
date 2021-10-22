package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pagarme/marshals/labs/vicki/desafioGo/logger"
	"github.com/pagarme/marshals/labs/vicki/desafioGo/prometheus"
)

var (
	Client httpClient
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	Client = &http.Client{Timeout: 30 * time.Second}
}

func insertAPIKey(apiKey string, paymentMethod []byte) []byte {
	var m map[string]interface{}
	json.Unmarshal(paymentMethod, &m)
	m["api_key"] = apiKey

	transactionWithAPIKey, err := json.Marshal(m)
	if err != nil {
		log.Println("Unable to insert the API key into the file")
	}
	return transactionWithAPIKey
}

func sendRequest(req *http.Request) (*http.Response, error) {
	res, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending the request: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		logger.ErrorLog.Printf("res.statuscode: %v, %s", res.StatusCode, res.Request.Method)
	}
	return res, err
}

func createTxRequest(endpoint, apiKey, authType string, paymentMethod []byte) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(paymentMethod))
	if err != nil {
		return nil, fmt.Errorf("error creating the request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	switch authType {
	case "inBody":
		paymentMethod = insertAPIKey(apiKey, paymentMethod)
		body := bytes.NewReader(paymentMethod)
		req.Body = ioutil.NopCloser(body)
		req.ContentLength = int64(body.Len())

	case "pathParam":
		param := req.URL.Query()
		param.Set("api_key", apiKey)
		req.URL.RawQuery = param.Encode()

	case "basicAuth":
		req.SetBasicAuth(apiKey, "x")

	default:
		return nil, fmt.Errorf("invalid authentication type: %s", authType)
	}
	return req, nil
}

// SendTx sends the API key using three types of authentication:
// 1. In the body of the request
// 2. Using query param
// 3. Using Basic Auth
func SendTx(endpoint, apiKey, authType string, paymentMethod []byte) ([]byte, error) {
	req, err := createTxRequest(endpoint, apiKey, authType, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("transaction with authentication %s failed: %s", authType, err.Error())
	}
	logger.InfoLog.Printf("json.request: %s %s, auth.type: %s", req.Method, req.URL.Path, authType)

	res, err := sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error calling sendRequest: %s", err.Error())
	}
	defer res.Body.Close()

	prometheus.RecordMetrics(authType, req.Method, res.Status)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}
	return body, nil
}
