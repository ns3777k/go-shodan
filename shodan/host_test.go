package shodan

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetHostsForQuery_DifferentVersionFormats(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(hostSearchPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "host/version"))
	})

	options := &HostQueryOptions{Query: "argentina"}
	_, err := client.GetHostsForQuery(options)

	assert.Nil(t, err)
}
