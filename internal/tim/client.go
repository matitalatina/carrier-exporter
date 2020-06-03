package tim

import (
	"net/http"
	"net/url"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Auth struct {
}

type Credentials struct {
	username string
	password string
}

type Client struct {
	Client HttpClient
}

type AuthorizeResponse struct {
	SessionID   string
	SessionData string
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

func getDefaultHeaders() http.Header {
	return http.Header{
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Accept-Language": {"it-IT,it;q=0.9,en-US;q=0.8,en;q=0.7"},
	}
}
