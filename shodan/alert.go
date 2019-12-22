package shodan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const (
	alertsInfoListPath = "/shodan/alert/info"
	alertInfoPath      = "/shodan/alert/%s/info"
	alertDeletePath    = "/shodan/alert/%s"
	alertCreatePath    = "/shodan/alert"
	alertNotifier      = "/shodan/alert/%s/notifier/%s"
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
	Notifiers  []*Notifier   `json:"notifiers"`
	// not documented for now :-(
	Triggers map[string]interface{} `json:"triggers"`
}

type alertCreateRequest struct {
	Name    string        `json:"name"`
	Expires int           `json:"expires"`
	Filters *AlertFilters `json:"filters"`
}

// CreateAlert creates a network alert for a defined IP/netblock which can be used to
// subscribe to changes/events that are discovered within that range.
func (c *Client) CreateAlert(ctx context.Context, name string, ip []string, expires int) (*Alert, error) {
	var alert Alert
	payload := &alertCreateRequest{
		Name:    name,
		Expires: expires,
		Filters: &AlertFilters{IP: ip},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("POST", alertCreatePath, nil, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &alert); err != nil {
		return nil, err
	}

	return &alert, nil
}

// GetAlerts returns a listing of all the network alerts that are currently active on the account.
func (c *Client) GetAlerts(ctx context.Context) ([]*Alert, error) {
	alerts := make([]*Alert, 0)

	req, err := c.NewRequest("GET", alertsInfoListPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &alerts); err != nil {
		return nil, err
	}

	return alerts, err
}

// GetAlert returns the information about a specific network alert.
func (c *Client) GetAlert(ctx context.Context, id string) (*Alert, error) {
	var alert Alert
	path := fmt.Sprintf(alertInfoPath, id)

	req, err := c.NewRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &alert); err != nil {
		return nil, err
	}

	return &alert, err
}

// DeleteAlert removes the specified network alert.
func (c *Client) DeleteAlert(ctx context.Context, id string) (bool, error) {
	path := fmt.Sprintf(alertDeletePath, id)

	req, err := c.NewRequest("DELETE", path, nil, nil)
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, nil); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) toggleAlertNotifier(
	ctx context.Context,
	method string,
	alertID string,
	notifierID string,
) (bool, error) {
	var response genericSuccessResponse
	path := fmt.Sprintf(alertNotifier, alertID, notifierID)

	req, err := c.NewRequest(method, path, nil, nil)
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &response); err != nil {
		return false, err
	}

	return response.Success, nil
}

// AddAlertNotifier adds the specified notifier to the network alert. Notifications are only sent if triggers have
// also been enabled.
func (c *Client) AddAlertNotifier(ctx context.Context, alertID string, notifierID string) (bool, error) {
	return c.toggleAlertNotifier(ctx, "PUT", alertID, notifierID)
}

// DeleteAlertNotifier removes the notification service from the alert. Notifications are only sent if triggers have
// also been enabled.
func (c *Client) DeleteAlertNotifier(ctx context.Context, alertID string, notifierID string) (bool, error) {
	return c.toggleAlertNotifier(ctx, "DELETE", alertID, notifierID)
}
