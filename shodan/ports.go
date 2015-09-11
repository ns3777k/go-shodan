package shodan

const (
	portsPath = "/shodan/ports"
)

type Port int

// GetPorts returns a list of port numbers that the crawlers are looking for
func (c *Client) GetPorts() ([]Port, error) {
	url, err := c.buildBaseURL(portsPath, nil)
	if err != nil {
		return nil, err
	}

	var ports []Port
	err = c.executeRequest("GET", url, &ports)

	return ports, err
}
