package shodan

const (
	protocolsPath = "/shodan/protocols"
)

// GetProtocols returns an object containing all the protocols that can be used when launching an Internet scan
func (c *Client) GetProtocols() (map[string]string, error) {
	url, err := c.buildBaseURL(protocolsPath, nil)
	if err != nil {
		return nil, err
	}

	var protocols map[string]string
	err = c.executeRequest("GET", url, &protocols, nil)

	return protocols, err
}
