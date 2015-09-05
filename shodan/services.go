package shodan

const (
	servicesPath = "/shodan/services"
)

type ServicePort        string
type ServiceDescription string
type ServiceCollection  map[ServicePort]ServiceDescription

func (c *Client) GetServices() (ServiceCollection, error) {
	url, err := c.buildUrl(servicesPath, nil)
	if err != nil {
		return nil, err
	}

	var services ServiceCollection
	err = c.executeRequest("GET", url, &services)

	return services, err
}
