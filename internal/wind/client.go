package wind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Credentials struct {
	Username string
	Password string
}

// Client retrieves Wind carrier data
type Client struct {
	Client      HttpClient
	authCookies []*http.Cookie
}

// Login makes login into Wind account
// it returns cookies to make authenticated calls
func (w *Client) Login(credentials Credentials) (string, error) {
	body := map[string]string{
		"username":   credentials.Username,
		"password":   credentials.Password,
		"rememberMe": "true",
	}
	jsonBody, err := json.Marshal(body)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://apigw.windtre.it/api/v4/login/credentials",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", err
	}

	req.Header = getDefaultHeaders()

	resp, err := w.Client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return resp.Header.Get("X-W3-Token"), nil
}

func (w *Client) GetStats(token string, lineId string, contractId string) (*Stats, error) {
	url := fmt.Sprintf("https://apigw.windtre.it/api/ob/v2/contract/lineunfolded?contractId=%s&lineId=%s&paymentType=POST", contractId, lineId)

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := w.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var stats Stats
	err = json.Unmarshal(body, &stats)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func getDefaultHeaders() http.Header {
	return http.Header{
		"X-Wind-Client":   {"app-and"},
		"X-Wind-Version":  {"ANDROID_V8.11.1"},
		"X-Brand":         {"ONEBRAND"},
		"X-W3-OS":         {"11"},
		"X-W3-Device":     {"Samsung SM-G970F"},
		"X-Language":      {"it"},
		"X-API-Client-Id": {"55527905-1fae-4f02-b7df-9c9a87749f69"},
		"Content-Type":    {"application/json; charset=UTF-8"},
		"Host":            {"apigw.windtre.it"},
	}
}
