package shodan

const (
	infoPath = "/api-info"
)

// APIInfo holds API information.
type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	HTTPS        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
	UnlockedLeft int    `json:"unlocked_left"`
}

// GetAPIInfo returns information about the API plan belonging to the given API key.
func (c *Client) GetAPIInfo() (*APIInfo, error) {
	url := c.buildBaseURL(infoPath, nil)

	var apiInfo APIInfo
	err := c.executeRequest("GET", url, &apiInfo, nil)

	return &apiInfo, err
}
