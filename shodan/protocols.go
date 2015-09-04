package shodan

const (
	protocolsPath = "/shodan/protocols"
)

type Protocol            string
type ProtocolDescription string
type ProtocolCollection  map[Protocol]ProtocolDescription

func (c *Client) GetProtocols() (*ProtocolCollection, error) {
	url, err := c.buildUrl(protocolsPath, nil)
	if err != nil {
		return nil, err
	}

	var protocols ProtocolCollection
	err = c.executeRequest("GET", url, &protocols)

	return &protocols, err
}
