package wind

import "net/http"

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
			Client: c.getHttpClient(),
		}
	}
	return c.Client
}

func (c *Container) GetService() *Service {
	return &Service{
		Client: *c.GetClient(),
	}
}
