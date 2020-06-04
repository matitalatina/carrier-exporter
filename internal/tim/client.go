package tim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Auth struct {
}

type Credentials struct {
	Username string
	Password string
}

type Client struct {
	Client HttpClient
}

type AuthorizeResponse struct {
	SessionID   string
	SessionData string
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
	IDTokenType  string `json:"id_token_type"`
}

type AuthApi struct {
	SessionID   string
	AccessToken string
}

type AggregateResponse struct {
	Status              string `json:"status"`
	Type                string `json:"type"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	ScadenzaGG          string `json:"scadenzaGG"`
	SimDeactivationDate string `json:"simDeactivationDate"`
	Action              struct {
		Path  string `json:"path"`
		Label string `json:"label"`
	} `json:"action"`
	Aggregateoffers []struct {
		Type       string `json:"type"`
		Value      string `json:"value"`
		Percent    int    `json:"percent"`
		MinPercent int    `json:"min_percent"`
		Unlimited  bool   `json:"unlimited"`
		ImageLink  string `json:"imageLink"`
		Total      string `json:"total,omitempty"`
		Measure    string `json:"measure,omitempty"`
	} `json:"aggregateoffers"`
}

func (c *Client) Authorize() (*AuthorizeResponse, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.tim.it/auth/oauth/v2/authorize?response_type=code&client_id=017da076-4447-4b82-89df-10cf1c15471a&redirect_uri=https%3A%2F%2Fmytim.tim.it%2Flogin-token.html&scope=openid+oob&prompt=none+SSO&nonce=cUJoIqwLIU3x&state=https%3A%2F%2Fmytim.tim.it%2F&code_challenge=YQLtzlD5BhUC_mt5JV73eQR6qqJW8JmNLrlMt5lfYzI&code_challenge_method=S256",
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

	redirectURL, err := url.Parse(resp.Header.Get("Location"))

	if err != nil {
		return nil, err
	}

	return &AuthorizeResponse{
		SessionID:   redirectURL.Query().Get("sessionID"),
		SessionData: redirectURL.Query().Get("sessionData"),
	}, nil
}

func (c *Client) Consent(auth AuthorizeResponse) (*string, error) {
	data := url.Values{
		"sessionID":     {auth.SessionID},
		"sessionData":   {auth.SessionData},
		"action":        {"grant"},
		"response_type": {"code"},
		"response_mode": {"query"},
	}
	req, err := http.NewRequest(
		"POST",
		"https://api.tim.it/auth/oauth/v2/authorize/consent",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	redirectURL, err := url.Parse(resp.Header.Get("Location"))

	if err != nil {
		return nil, err
	}

	code := redirectURL.Query().Get("code")
	return &code, nil
}

func (c *Client) Token(code string) (*AuthTokens, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"https://mytim.tim.it/login-token.html"},
		"client_id":     {"017da076-4447-4b82-89df-10cf1c15471a"},
		"code_verifier": {"2EaYljO0DwRiQe4SK6GCDJ72iDUWvLrvtSM6lNCf45ndDMMwa1OgZyX3i4zLnM5kFbAr2hEpDdS3Ye8IWvf83wV3u4ba9E6mDyjSehdS7VXQ6KSR5HyCUuRii3a1JAHF"},
	}
	req, err := http.NewRequest(
		"POST",
		"https://api.tim.it/auth/oauth/v2/token",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var authTokens AuthTokens
	err = json.Unmarshal(body, &authTokens)

	if err != nil {
		return nil, err
	}

	return &authTokens, nil
}

func (c *Client) Aggregate(authApi AuthApi, phone string) (*AggregateResponse, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.tim.it/api/consistenze/%s/offerte/aggregato", phone),
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authApi.AccessToken)
	req.Header.Set("BusinessID", "141e639d4745401f8051850b")
	req.Header.Set("Channel", "MYTIMWEB")
	req.Header.Set("TransactionID", "F267C2CC5B504B9F985DBDF2")
	req.Header.Set("SessionID", authApi.SessionID)

	resp, err := c.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var aggregateResponse AggregateResponse
	err = json.Unmarshal(body, &aggregateResponse)

	if err != nil {
		return nil, err
	}

	return &aggregateResponse, nil
}

func (c *Client) Login(auth AuthorizeResponse, credentials Credentials) (*AuthorizeResponse, error) {
	data := url.Values{
		"sessionID":     {auth.SessionID},
		"sessionData":   {auth.SessionData},
		"username":      {credentials.Username},
		"password":      {credentials.Password},
		"remember":      {"true"},
		"response_type": {""},
		"response_mode": {""},
		"action":        {"login"},
	}
	req, err := http.NewRequest(
		"POST",
		"https://api.tim.it/auth/oauth/v2/authorize/login",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header = getDefaultHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	sessionID, _ := doc.Find("form input[name=sessionID]").First().Attr("value")
	sessionData, _ := doc.Find("form input[name=sessionData]").First().Attr("value")

	if err != nil {
		return nil, err
	}

	return &AuthorizeResponse{
		SessionID:   sessionID,
		SessionData: sessionData,
	}, nil
}

func getDefaultHeaders() http.Header {
	return http.Header{
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Accept-Language": {"it-IT,it;q=0.9,en-US;q=0.8,en;q=0.7"},
	}
}
