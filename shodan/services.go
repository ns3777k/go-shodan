package shodan

const (
	servicesPath = "/shodan/services"
)

// GetServices returns an object containing all the services that the Shodan crawlers look at
// It can also be used as a quick and practical way to resolve a port number to the name of a service
func (c *Client) GetServices() (map[string]string, error) {
	url := c.buildBaseURL(servicesPath, nil)

	var services map[string]string
	err := c.executeRequest("GET", url, &services, nil)

	return services, err
}
