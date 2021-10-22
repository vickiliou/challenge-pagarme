package transaction

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

var (
	a = struct {
		apiKey, endpoint string
		authType         [3]string
		paymentMethod    []byte
		resBody          [3]string
	}{
		apiKey:   "fakeApiKey",
		endpoint: "https://fakeEndpoint",
		authType: [3]string{
			"inBody",
			"pathParam",
			"basicAuth",
		},
		paymentMethod: []byte(`{"fakePaymentMethod":"a"}`),
		resBody: [3]string{
			`{"api_key":"fakeApiKey","fakePaymentMethod":"a"}`,
			`{"fakePaymentMethod":"a"}`,
			`{"fakePaymentMethod":"a"}`,
		},
	}
)

type MockClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestInsertAPIKey(t *testing.T) {
	expect := a.resBody[0]
	actual := insertAPIKey(a.apiKey, a.paymentMethod)
	if string(actual) != expect {
		t.Errorf("Error insertAPIKey, got: %s, want: %s", actual, expect)
	}
}

func TestSendRequest(t *testing.T) {
	expect := http.StatusOK
	r := ioutil.NopCloser(bytes.NewReader(a.paymentMethod))
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}
	req, err := createTxRequest(a.endpoint, a.apiKey, a.authType[0], a.paymentMethod)
	if err != nil {
		t.Errorf("Error createTxRequest: %s", err.Error())
	}
	actual, err := sendRequest(req)
	if err != nil {
		t.Errorf("Error sendRequest: %s", err.Error())
		return
	}
	if actual.StatusCode != expect {
		t.Errorf("Error TestSendRequest, got: %v, want: %v", actual.StatusCode, expect)
		return
	}
}

func TestSendRequestFail(t *testing.T) {
	expect := "error while sending the request: Mock Error"
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       nil,
			}, errors.New("Mock Error")
		},
	}
	req, err := createTxRequest(a.endpoint, a.apiKey, a.authType[0], a.paymentMethod)
	if err != nil {
		t.Errorf("Error createTxRequest: %s", err.Error())
	}
	_, err = sendRequest(req)
	if err.Error() != expect {
		t.Errorf("Error sendRequest: %s", err.Error())
	}
}

func TestCreateTxRequest(t *testing.T) {
	for i, dataAuth := range a.authType {
		u, _ := url.Parse(a.endpoint)
		if dataAuth == "pathParam" {
			ua, _ := url.Parse(a.endpoint)
			u, _ = ua.Parse("?api_key=" + a.apiKey)
		}
		req, err := createTxRequest(a.endpoint, a.apiKey, dataAuth, a.paymentMethod)
		if err != nil {
			t.Errorf("Error createTxRequest: %s", err.Error())
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Error header, got: %s, want: %s", req.Header.Get("Content-Type"), "application/json")
		}
		if !reflect.DeepEqual(req.URL, u) {
			if !reflect.DeepEqual(req.URL, u) {
				t.Errorf("Error endpoint %s, got: %v, want: %v", dataAuth, req.URL, u)
			}
		}
		body, _ := ioutil.ReadAll(req.Body)
		if string(body) != a.resBody[i] {
			t.Errorf("Error body, got: %s, want: %s", string(body), a.resBody[i])
		}
	}
}

func TestCreateTxRequestWithInvalidAuth(t *testing.T) {
	expectedWrongAuth := "invalid authentication type: fakeAuth"
	_, err := createTxRequest(a.endpoint, a.apiKey, "fakeAuth", a.paymentMethod)
	if err.Error() != expectedWrongAuth {
		t.Errorf("Error createTxRequest, got: %s, want: %s", err.Error(), expectedWrongAuth)
	}
}

func TestSendTx(t *testing.T) {
	expect := string(a.paymentMethod)
	r := ioutil.NopCloser(bytes.NewReader(a.paymentMethod))
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}
	actual, err := SendTx(a.endpoint, a.apiKey, a.authType[0], a.paymentMethod)
	if err != nil {
		t.Errorf("Error SendTx failed, %s", err.Error())
		return
	}
	if string(actual) != expect {
		t.Errorf("Error TestSendTx, got: %s, want: %s", string(actual), expect)
		return
	}
}

func TestSendTxFail(t *testing.T) {
	expect := "error calling sendRequest: error while sending the request: Mock Error"
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       nil,
			}, errors.New("Mock Error")
		},
	}
	_, err := SendTx(a.endpoint, a.apiKey, a.authType[0], a.paymentMethod)
	if err.Error() != expect {
		t.Errorf("Error SendTx, got: %s, want: %s", err.Error(), expect)
		return
	}
}
