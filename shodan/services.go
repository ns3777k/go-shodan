package shodan

const (
	servicesPath = "/shodan/services"
)

type ServiceCollection map[string]string

// GetServices returns an object containing all the services that the Shodan crawlers look at
// It can also be used as a quick and practical way to resolve a port number to the name of a service
func (c *Client) GetServices() (ServiceCollection, error) {
	url, err := c.buildBaseURL(servicesPath, nil)
	if err != nil {
		return nil, err
	}

	var services ServiceCollection
	err = c.executeRequest("GET", url, &services, nil)

	return services, err
}
