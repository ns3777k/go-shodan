package shodan

const (
	protocolsPath = "/shodan/protocols"
)

type ProtocolCollection map[string]string

// GetProtocols returns an object containing all the protocols that can be used when launching an Internet scan
func (c *Client) GetProtocols() (ProtocolCollection, error) {
	url, err := c.buildBaseURL(protocolsPath, nil)
	if err != nil {
		return nil, err
	}

	var protocols ProtocolCollection
	err = c.executeRequest("GET", url, &protocols)

	return protocols, err
}
