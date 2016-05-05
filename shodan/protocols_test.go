package shodan

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetProtocols(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(protocolsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "protocols"))
	})

	protocolsExpected := map[string]string{
		"andromouse": "Checks whether the device is running the remote mouse AndroMouse service.",
		"zookeeper":  "Grab statistical information from a Zookeeper node",
	}
	protocols, err := client.GetProtocols()

	assert.Nil(t, err)
	assert.Len(t, protocols, len(protocolsExpected))
	assert.EqualValues(t, protocolsExpected, protocols)
}

func TestClient_GetProtocols_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetProtocols()
	assert.NotNil(t, err)
}
