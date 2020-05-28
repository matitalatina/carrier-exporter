package vodafone

import (
	"net/http"
	"time"
)

type Container struct {
	Client     *Client
	httpClient *http.Client
}

func (c *Container) getHttpClient() *http.Client {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c.httpClient
}

func (c *Container) GetClient() *Client {
	if c.Client == nil {
		c.Client = &Client{
			Client: &http.Client{
				Timeout: time.Second * 10,
			},
		}
	}
	return c.Client
}
