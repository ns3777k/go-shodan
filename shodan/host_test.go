package shodan

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetServicesForHost_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetServicesForHost("192.168.1.1", new(HostServicesOptions))
	assert.NotNil(t, err)
}

func TestClient_GetHostsCountForQuery_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetHostsCountForQuery(new(HostQueryOptions))
	assert.NotNil(t, err)
}

func TestClient_GetHostsForQuery_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetHostsForQuery(new(HostQueryOptions))
	assert.NotNil(t, err)
}

func TestClient_BreakQueryIntoTokens_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.BreakQueryIntoTokens("anything here")
	assert.NotNil(t, err)
}
