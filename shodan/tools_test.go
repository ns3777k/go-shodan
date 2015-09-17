package shodan

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetMyIP(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	testIP := "192.168.22.34"

	mux.HandleFunc(ipPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, strconv.Quote(testIP))
	})

	ip, err := client.GetMyIP()

	assert.Nil(t, err)
	assert.Equal(t, testIP, ip)
}
