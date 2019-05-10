package shodan

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetPorts(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(portsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "ports")) //nolint:errcheck
	})

	portsExpected := []int{22, 771, 5353, 110, 8139}
	ports, err := client.GetPorts(context.TODO())

	assert.Nil(t, err)
	assert.Len(t, ports, len(portsExpected))
	assert.EqualValues(t, portsExpected, ports)
}
