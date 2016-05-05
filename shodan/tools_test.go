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

func TestClient_GetMyIP_invalidBaseURL(t *testing.T) {
	client := NewClient(nil, testClientToken)
	client.BaseURL = ":/1232.22"
	_, err := client.GetMyIP()
	assert.NotNil(t, err)
}
