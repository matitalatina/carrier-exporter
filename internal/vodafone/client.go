package vodafone

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

type Client struct {
	Client HttpClient
}

type ApiAuth struct {
	BwbSessionId string
	Token        string
}
type GuestCredentials struct {
	BwbSessionId string
	Response     GuestResponse
}

func (c *Client) GetToken() (*GuestCredentials, error) {
	req, err := http.NewRequest(
		"POST",
		"https://my190.vodafone.it/api/v3/auth/guest",
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()

	resp, err := c.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var guestResponse GuestResponse
	err = json.Unmarshal(body, &guestResponse)

	if err != nil {
		return nil, err
	}

	return &GuestCredentials{
		BwbSessionId: resp.Header.Get("X-Bwb-SessionId"),
		Response:     guestResponse,
	}, nil
}

func (c *Client) Login(auth ApiAuth, credentials Credentials) (*LoginResponse, error) {
	bodyReq := map[string]string{
		"username": credentials.Username,
		"password": credentials.Password,
	}
	jsonBody, err := json.Marshal(bodyReq)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://my190.vodafone.it/api/v3/auth/login?includeItem=sim&includeItem=landline",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()
	req.Header.Set("X-Bwb-SessionId", auth.BwbSessionId)
	req.Header.Set("X-Auth-Token", auth.Token)

	resp, err := c.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)

	if err != nil {
		return nil, err
	}

	return &loginResponse, nil
}

func (c *Client) GetCounters(auth ApiAuth, phone string) (*CountersResponse, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://my190.vodafone.it/api/v3/sim/%s/counters", phone),
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()
	req.Header.Set("X-Bwb-SessionId", auth.BwbSessionId)
	req.Header.Set("X-Auth-Token", auth.Token)

	resp, err := c.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var counterResponse CountersResponse
	err = json.Unmarshal(body, &counterResponse)

	if err != nil {
		return nil, err
	}

	return &counterResponse, nil
}

func getDefaultHeaders() http.Header {
	return http.Header{
		"X-Bwb-AppId": {"MyVodafone"},
		// These ids are random (32 chars)
		// https://passwordsgenerator.net/
		"X-Bwb-InstallationId": {"i0uifrsnao7zes0k1xkom0wg43iqsa5h"},
		"X-Device-Id":          {"i0uifrsnao7zes0k1xkom0wg43iqsa5h"},
		"clientId":             {"4002"},
		"X-Bwb-User-Agent":     {"MyVodafone/android/s/11/12.12.2/2.625/store/\"samsung\"/\"SM-G970F\""},
		"X-Device-UserAgent":   {"MyVodafone/android/s/1/12.12.2/2.625/store/\"samsung\"/\"SM-G970F\""},
		"User-Agent":           {"Dalvik/2.1.0 (Linux; U; Android 11; SM-G970F Build/QP1A.190711.020)"},
		"X-Bwb-AppVersion":     {"12.12.2"},
		"Content-Type":         {"application/json; charset=UTF-8"},
		"X-Network-Info":       {"w/nd/0.71/222"},
	}
}
