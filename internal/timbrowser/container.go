package timbrowser

type Container struct {
	Service *Service
}

func (c *Container) GetService() *Service {
	if c.Service == nil {
		c.Service = &Service{}
	}
	return c.Service
}
