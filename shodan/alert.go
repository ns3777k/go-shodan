package shodan

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	alertsInfoListPath = "/shodan/alert/info"
	alertInfoPath      = "/shodan/alert/%s/info"
	alertDeletePath    = "/shodan/alert/%s"
	alertCreatePath    = "/shodan/alert"
)

type AlertFilters struct {
	IP []string `json:"ip"`
}

type Alert struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Created    string        `json:"created"`
	Expiration string        `json:"expiration"`
	Expires    int           `json:"expires"`
	Expired    bool          `json:"expired"`
	Size       int           `json:"size"`
	Filters    *AlertFilters `json:"filters"`
}

type alertCreateRequest struct {
	Name    string        `json:"name"`
	Expires int           `json:"expires"`
	Filters *AlertFilters `json:"filters"`
}

// Use this method to create a network alert for a defined IP/ netblock which
// can be used to subscribe to changes/ events that are discovered within that range.
func (c *Client) CreateAlert(name string, ip []string, expires int) (*Alert, error) {
	url, err := c.buildBaseURL(alertCreatePath, nil)
	if err != nil {
		return nil, err
	}

	payload := &alertCreateRequest{
		Name:    name,
		Expires: expires,
		Filters: &AlertFilters{
			IP: ip,
		},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var alert Alert
	err = c.executeRequest("POST", url, &alert, bytes.NewReader(b))

	return &alert, err
}

// GetAlerts returns a listing of all the network alerts
// that are currently active on the account.
func (c *Client) GetAlerts() ([]*Alert, error) {
	url, err := c.buildBaseURL(alertsInfoListPath, nil)
	if err != nil {
		return nil, err
	}

	alerts := make([]*Alert, 0, 0)
	err = c.executeRequest("GET", url, &alerts, nil)

	return alerts, err
}

// GetAlert returns the information about a specific network alert.
func (c *Client) GetAlert(id string) (*Alert, error) {
	path := fmt.Sprintf(alertInfoPath, id)
	url, err := c.buildBaseURL(path, nil)
	if err != nil {
		return nil, err
	}

	var alert Alert
	err = c.executeRequest("GET", url, &alert, nil)

	return &alert, err
}

// DeleteAlert removes the specified network alert.
func (c *Client) DeleteAlert(id string) (bool, error) {
	path := fmt.Sprintf(alertDeletePath, id)
	url, err := c.buildBaseURL(path, nil)
	if err != nil {
		return false, err
	}

	err = c.executeRequest("DELETE", url, nil, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}
