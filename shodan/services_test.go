package shodan

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetServices(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(servicesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "services"))
	})

	servicesExpected := map[string]string{
		"623":  "IPMI",
		"8181": "GlassFish Server",
		"53":   "DNS",
	}
	services, err := client.GetServices()

	assert.Nil(t, err)
	assert.Len(t, services, len(servicesExpected))
	assert.EqualValues(t, servicesExpected, services)
}

func TestClient_GetServices_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetServices()
	assert.NotNil(t, err)
}
