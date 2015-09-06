package shodan

const (
	profilePath = "/account/profile"
)

type Profile struct {
	Member  bool   `json:"member"`
	Credits int    `json:"credits"`
	Name    string `json:"display_name"`
	Created string `json:"created"`
}

// GetAccountProfile returns information about the Shodan account linked to the API key
func (c *Client) GetAccountProfile() (*Profile, error) {
	url, err := c.buildURL(profilePath, nil)
	if err != nil {
		return nil, err
	}

	var profile Profile
	err = c.executeRequest("GET", url, &profile)

	return &profile, err
}
