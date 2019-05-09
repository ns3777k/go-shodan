package shodan

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetProtocols(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(protocolsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "protocols")) //nolint:errcheck
	})

	protocolsExpected := map[string]string{
		"andromouse": "Checks whether the device is running the remote mouse AndroMouse service.",
		"zookeeper":  "Grab statistical information from a Zookeeper node",
	}
	protocols, err := client.GetProtocols(context.TODO())

	assert.Nil(t, err)
	assert.Len(t, protocols, len(protocolsExpected))
	assert.EqualValues(t, protocolsExpected, protocols)
}
