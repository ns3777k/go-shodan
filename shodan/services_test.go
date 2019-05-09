package shodan

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetServices(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(servicesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "services")) //nolint:errcheck
	})

	servicesExpected := map[string]string{
		"623":  "IPMI",
		"8181": "GlassFish Server",
		"53":   "DNS",
	}

	services, err := client.GetServices(context.TODO())

	assert.Nil(t, err)
	assert.Len(t, services, len(servicesExpected))
	assert.EqualValues(t, servicesExpected, services)
}
