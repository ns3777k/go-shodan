package shodan

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetPorts(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(portsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "ports"))
	})

	portsExpected := []int{22, 771, 5353, 110, 8139}
	ports, err := client.GetPorts()

	assert.Nil(t, err)
	assert.Len(t, ports, len(portsExpected))
	assert.EqualValues(t, portsExpected, ports)
}

func TestClient_GetPorts_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetPorts()
	assert.NotNil(t, err)
}
