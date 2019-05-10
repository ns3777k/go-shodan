package shodan

import (
	"context"
	"fmt"
)

const (
	alertTriggersListPath     = "/shodan/alert/triggers"
	alertTriggerEnablePath    = "/shodan/alert/%s/trigger/%s"
	alertTriggerWhitelistPath = "/shodan/alert/%s/trigger/%s/ignore/%s"
)

// AlertTrigger represents a trigger.
type AlertTrigger struct {
	Name        string `json:"name"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

type AlertTriggerIdent struct {
	AlertID     string
	TriggerName string
}

type AlertTriggerServiceIdent struct {
	*AlertTriggerIdent
	ServiceName string
}

type alertTriggerSuccessResponse struct {
	Success bool `json:"success"`
}

// Returns a list of all the triggers that can be enabled on network alerts.
func (c *Client) GetAlertTriggers(ctx context.Context) ([]*AlertTrigger, error) {
	triggers := make([]*AlertTrigger, 0)

	req, err := c.NewRequest("GET", alertTriggersListPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &triggers); err != nil {
		return nil, err
	}

	return triggers, err
}

func (c *Client) toggleAlertTrigger(ctx context.Context, method string, trigger *AlertTriggerIdent) (bool, error) {
	var response alertTriggerSuccessResponse
	path := fmt.Sprintf(alertTriggerEnablePath, trigger.AlertID, trigger.TriggerName)

	req, err := c.NewRequest(method, path, nil, nil)
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &response); err != nil {
		return false, err
	}

	return response.Success, nil
}

func (c *Client) toggleAlertTriggerWhitelist(
	ctx context.Context,
	method string,
	service *AlertTriggerServiceIdent,
) (bool, error) {
	var response alertTriggerSuccessResponse
	path := fmt.Sprintf(alertTriggerWhitelistPath, service.AlertID, service.TriggerName, service.ServiceName)

	req, err := c.NewRequest(method, path, nil, nil)
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &response); err != nil {
		return false, err
	}

	return response.Success, nil
}

// Get notifications when the specified trigger is met.
func (c *Client) EnableAlertTrigger(ctx context.Context, ident *AlertTriggerIdent) (bool, error) {
	return c.toggleAlertTrigger(ctx, "PUT", ident)
}

// Stop getting notifications for the specified trigger.
func (c *Client) DisableAlertTrigger(ctx context.Context, ident *AlertTriggerIdent) (bool, error) {
	return c.toggleAlertTrigger(ctx, "DELETE", ident)
}

// Ignore the specified service when it is matched for the trigger.
func (c *Client) AddServiceToAlertTriggerWhitelist(
	ctx context.Context,
	service *AlertTriggerServiceIdent,
) (bool, error) {
	return c.toggleAlertTriggerWhitelist(ctx, "PUT", service)
}

// Start getting notifications again for the specified trigger.
func (c *Client) RemoveServiceFromAlertTriggerWhitelist(
	ctx context.Context,
	service *AlertTriggerServiceIdent,
) (bool, error) {
	return c.toggleAlertTriggerWhitelist(ctx, "DELETE", service)
}
