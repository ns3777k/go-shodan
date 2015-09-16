package shodan

import (
	"testing"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetServices(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(servicesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "services"))
	})

	services, err := client.GetServices()

	assert.Nil(t, err);
	assert.Len(t, services, 155)
}
