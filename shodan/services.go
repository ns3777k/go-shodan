package shodan

import (
	"context"
)

const (
	servicesPath = "/shodan/services"
)

// GetServices returns an object containing all the services that the Shodan crawlers look at
// It can also be used as a quick and practical way to resolve a port number to the name of a service
func (c *Client) GetServices(ctx context.Context) (map[string]string, error) {
	var services map[string]string

	req, err := c.NewRequest("GET", servicesPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &services); err != nil {
		return nil, err
	}

	return services, err
}
