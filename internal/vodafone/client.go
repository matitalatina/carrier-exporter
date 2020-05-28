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

type GuestCredentials struct {
	BwbSessionid string
	Response     GuestResponse
}

type GuestResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Result      struct {
		Token                  string `json:"token"`
		WidgetAutoRefreshDelay string `json:"widgetAutoRefreshDelay"`
		WidgetNextRefreshDelay string `json:"widgetNextRefreshDelay"`
	} `json:"result"`
}

type LoginResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Result      struct {
		Customer struct {
			Email      string `json:"email"`
			Firstname  string `json:"firstname"`
			FiscalCode string `json:"fiscalCode"`
			Lastname   string `json:"lastname"`
			TypeID     string `json:"typeId"`
			Username   string `json:"username"`
		} `json:"customer"`
		Items []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"items"`
		PartyIdentifier string `json:"partyIdentifier"`
		PsgToken        string `json:"psgToken"`
	} `json:"result"`
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
		BwbSessionid: resp.Header.Get("X-Bwb-SessionId"),
		Response:     guestResponse,
	}, nil
}

func (c *Client) Login(jwtToken string, bwbSessionId string, credentials Credentials) (*LoginResponse, error) {
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
	req.Header.Set("X-Bwb-SessionId", bwbSessionId)
	req.Header.Set("X-Auth-Token", jwtToken)

	resp, err := c.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)

	if err != nil {
		return nil, err
	}

	return &loginResponse, nil
}

func getDefaultHeaders() http.Header {
	return http.Header{
		"X-Bwb-AppId": {"MyVodafone"},
		// These ids are random (32 chars)
		// https://passwordsgenerator.net/
		"X-Bwb-InstallationId": {"r02evmxqhcmffg9yd3qcojerk38tuqer"},
		"X-Device-Id":          {"r02evmxqhcmffg9yd3qcojerk38tuqer"},
		"clientId":             {"4002"},
		"X-Bwb-User-Agent":     {"MyVodafone/android/s/10/12.4.0/2.625/store/\"samsung\"/\"SM-G970F\""},
		"X-Device-UserAgent":   {"MyVodafone/android/s/10/12.4.0/2.625/store/\"samsung\"/\"SM-G970F\""},
		"User-Agent":           {"Dalvik/2.1.0 (Linux; U; Android 10; SM-G970F Build/QP1A.190711.020)"},
		"X-Bwb-AppVersion":     {"12.4.0"},
		"Content-Type":         {"application/json; charset=UTF-8"},
		"X-Network-Info":       {"w/nd/0.71/222"},
	}
}
