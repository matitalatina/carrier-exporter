package vodafone

import (
	"net/http"
	"time"
)

type Container struct {
	Client     *Client
	Service    *Service
	httpClient *http.Client
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

func (c *Container) GetService() *Service {
	return &Service{
		Client: *c.GetClient(),
	}
}
