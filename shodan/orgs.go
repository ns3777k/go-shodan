package shodan

import (
	"context"
	"fmt"
)

const (
	organizationPath       = "/org"
	organizationMemberPath = "/org/member/%s"
)

// OrganizationMember is a company's employee.
type OrganizationMember struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Organization contains everything about a company.
type Organization struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Created     string                `json:"created"`
	Admins      []*OrganizationMember `json:"admins"`
	Members     []*OrganizationMember `json:"members"`
	UpgradeType string                `json:"upgrade_type"`
	Domains     []string              `json:"domains"`
	Logo        *string               `json:"logo"`
}

// AddMemberToOrganizationOptions is options for adding members.
type AddMemberToOrganizationOptions struct {
	Notify bool `url:"notify,omitempty"`
}

type organizationMemberSuccessResponse struct {
	Success bool `json:"success"`
}

// Get information about your organization such as the list of its members, upgrades, authorized domains and more.
func (c *Client) GetOrganization(ctx context.Context) (*Organization, error) {
	var organization Organization

	req, err := c.NewRequest("GET", organizationPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &organization); err != nil {
		return nil, err
	}

	return &organization, nil
}

// Add a Shodan user to the organization and upgrade them.
func (c *Client) AddMemberToOrganization(
	ctx context.Context,
	user string,
	options *AddMemberToOrganizationOptions,
) (bool, error) {
	var response organizationMemberSuccessResponse
	path := fmt.Sprintf(organizationMemberPath, user)

	req, err := c.NewRequest("PUT", path, options, nil)
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &response); err != nil {
		return false, err
	}

	return response.Success, nil
}

// Remove and downgrade the provided member from the organization.
func (c *Client) RemoveMemberFromOrganization(ctx context.Context, user string) (bool, error) {
	var response organizationMemberSuccessResponse
	path := fmt.Sprintf(organizationMemberPath, user)

	req, err := c.NewRequest("DELETE", path, nil, nil)
	if err != nil {
		return false, err
	}

	if err := c.Do(ctx, req, &response); err != nil {
		return false, err
	}

	return response.Success, nil
}
