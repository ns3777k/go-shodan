package shodan

import (
	"context"
	"net/url"
	"strings"
)

const (
	notifierPath         = "/notifier"
	notifierProviderPath = "/notifier/provider"
)

type Notifier struct {
	ID          string            `json:"id"`
	Provider    string            `json:"provider"`
	Description string            `json:"description"`
	Args        map[string]string `json:"args"`
}

type NotifierProvider struct {
	Required []string `json:"required"`
}

type notifierSuccessResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// GetNotifiers gets a list of all the notifiers that the user has created.
func (c *Client) GetNotifiers(ctx context.Context) ([]*Notifier, error) {
	matchesResponse := struct {
		Matches []*Notifier `json:"matches"`
	}{}

	req, err := c.NewRequest("GET", notifierPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &matchesResponse); err != nil {
		return nil, err
	}

	return matchesResponse.Matches, nil
}

// GetNotifierProviders gets a list of all the notification providers that are available and the parameters to submit
// when creating them.
func (c *Client) GetNotifierProviders(ctx context.Context) (map[string]*NotifierProvider, error) {
	providers := make(map[string]*NotifierProvider)

	req, err := c.NewRequest("GET", notifierProviderPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &providers); err != nil {
		return nil, err
	}

	return providers, nil
}

// GetNotifier gets information about a notifier.
func (c *Client) GetNotifier(ctx context.Context, id string) (*Notifier, error) {
	var notifier Notifier

	req, err := c.NewRequest("GET", notifierPath+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &notifier); err != nil {
		return nil, err
	}

	return &notifier, nil
}

// DeleteNotifier deletes the notification service created for the user.
func (c *Client) DeleteNotifier(ctx context.Context, id string) (bool, error) {
	var r genericSuccessResponse

	req, err := c.NewRequest("DELETE", notifierPath+"/"+id, nil, nil)
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &r); err != nil {
		return false, err
	}

	return r.Success, nil
}

// CreateNotifier creates a new notification service endpoint that Shodan services can send notifications through.
func (c *Client) CreateNotifier(ctx context.Context, notifier *Notifier) (bool, error) {
	var r notifierSuccessResponse

	body := url.Values{}
	body.Add("provider", notifier.Provider)
	body.Add("description", notifier.Description)

	for k, v := range notifier.Args {
		body.Add(k, v)
	}

	req, err := c.NewRequest("POST", notifierPath, nil, strings.NewReader(body.Encode()))
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &r); err != nil {
		return false, err
	}

	notifier.ID = r.ID

	return r.Success, nil
}

// UpdateNotifier updates the parameters of a notifier.
func (c *Client) UpdateNotifierArgs(ctx context.Context, id string, args map[string]string) (bool, error) {
	var r genericSuccessResponse

	body := url.Values{}

	for k, v := range args {
		body.Add(k, v)
	}

	req, err := c.NewRequest("PUT", notifierPath+"/"+id, nil, strings.NewReader(body.Encode()))
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &r); err != nil {
		return false, err
	}

	return r.Success, nil
}
