package shodan

const (
	portsPath = "/shodan/ports"
)

// GetPorts returns a list of port numbers that the crawlers are looking for
func (c *Client) GetPorts() ([]int, error) {
	url, err := c.buildBaseURL(portsPath, nil)
	if err != nil {
		return nil, err
	}

	var ports []int
	err = c.executeRequest("GET", url, &ports, nil)

	return ports, err
}
