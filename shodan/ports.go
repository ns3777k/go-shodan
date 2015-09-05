package shodan

const (
	portsPath = "/shodan/ports"
)

type Port int

func (c *Client) GetPorts() ([]Port, error) {
	url, err := c.buildUrl(portsPath, nil)
	if err != nil {
		return nil, err
	}

	var ports []Port
	err = c.executeRequest("GET", url, &ports)

	return ports, err
}
