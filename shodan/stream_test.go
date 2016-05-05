package shodan

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetBannersByPorts_invalidStreamBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.StreamBaseURL = ":/1232.22"
	err := client.GetBannersByPorts([]int{22, 45})
	assert.NotNil(t, err)
}

func TestClient_GetBanners_invalidStreamBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.StreamBaseURL = ":/1232.22"
	err := client.GetBanners()
	assert.NotNil(t, err)
}
