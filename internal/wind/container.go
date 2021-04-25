package wind

import (
	"crypto/tls"
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
		// Today 2021-04-25T08:25:45.133Z, Wind Apis have expired certificate :D
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.Client = &Client{
			Client: &http.Client{
				Timeout:   time.Second * 10,
				Transport: tr,
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
