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

// AlertFilters holds alert criteria (only ip for now).
type AlertFilters struct {
	IP []string `json:"ip"`
}

// Alert represents a trigger to react to network scan request.
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

// CreateAlert creates a network alert for a defined IP/ netblock which can be used to
// subscribe to changes/ events that are discovered within that range.
func (c *Client) CreateAlert(name string, ip []string, expires int) (*Alert, error) {
	url := c.buildBaseURL(alertCreatePath, nil)

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
	url := c.buildBaseURL(alertsInfoListPath, nil)

	alerts := make([]*Alert, 0, 0)
	err := c.executeRequest("GET", url, &alerts, nil)

	return alerts, err
}

// GetAlert returns the information about a specific network alert.
func (c *Client) GetAlert(id string) (*Alert, error) {
	path := fmt.Sprintf(alertInfoPath, id)
	url := c.buildBaseURL(path, nil)

	var alert Alert
	err := c.executeRequest("GET", url, &alert, nil)

	return &alert, err
}

// DeleteAlert removes the specified network alert.
func (c *Client) DeleteAlert(id string) (bool, error) {
	path := fmt.Sprintf(alertDeletePath, id)
	url := c.buildBaseURL(path, nil)

	err := c.executeRequest("DELETE", url, nil, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}
