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
func (w *Client) Login(credentials Credentials) ([]*http.Cookie, error) {
	body := map[string]string{
		"username": credentials.Username,
		"password": credentials.Password,
	}
	jsonBody, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://apigateway-selfcare.apps.windtre.it/api/v1/authentication/authenticate-app",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := w.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	cookies := resp.Cookies()

	return cookies, nil
}

func (w *Client) GetStats(authCookies []*http.Cookie, lineId string, contractId string) (*Stats, error) {
	url := fmt.Sprintf("https://apigateway-selfcare.apps.windtre.it/api/v1/contract/lineunfolded?lineId=%s&contractId=%s", lineId, contractId)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	for _, c := range authCookies {
		req.AddCookie(c)
	}

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
