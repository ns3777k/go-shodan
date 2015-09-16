package shodan

import (
	"testing"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetProtocols(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(protocolsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "protocols"))
	})

	protocols, err := client.GetProtocols()

	assert.Nil(t, err);
	assert.Len(t, protocols, 116)
}
