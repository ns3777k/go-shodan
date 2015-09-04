package shodan

const (
	infoPath = "/api-info"
)

type ApiInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	Https        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}

func (c *Client) GetApiInfo() (*ApiInfo, error) {
	url, err := c.buildUrl(infoPath, nil)
	if err != nil {
		return nil, err
	}

	var apiInfo ApiInfo
	err = c.executeRequest("GET", url, &apiInfo)

	return &apiInfo, err
}
