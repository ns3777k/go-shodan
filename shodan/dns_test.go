package shodan

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"net"
)

func TestClient_GetDNSResolve(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	expectedHostnames := []string{"google.com", "bing.com", "idonotexist.local"}

	mux.HandleFunc(resolvePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		hostnames := r.URL.Query().Get("hostnames")
		assert.NotEmpty(t, hostnames)

		splited := strings.Split(hostnames, ",")
		assert.Len(t, splited, len(expectedHostnames))

		w.Write(getStub(t, "dns_resolve"))
	})

	resolve, err := client.GetDNSResolve(expectedHostnames)

	assert.Nil(t, err)
	assert.Len(t, resolve, len(expectedHostnames))

	for _, host := range expectedHostnames {
		_, ok := resolve[host]
		assert.True(t, ok)
	}
}

func TestClient_GetDNSReverse(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	expectedIPs := []string{"74.125.227.244", "92.63.108.40"}

	mux.HandleFunc(reversePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		ips := r.URL.Query().Get("ips")
		assert.NotEmpty(t, ips)

		splited := strings.Split(ips, ",")
		assert.Len(t, splited, len(expectedIPs))

		w.Write(getStub(t, "dns_reverse"))
	})

	reversed, err := client.GetDNSReverse(expectedIPs)

	assert.Nil(t, err)
	assert.Len(t, reversed, len(expectedIPs))

	for _, ip := range expectedIPs {
		_, ok := reversed[ip]
		assert.True(t, ok)
	}
}

func TestClient_GetDNSReverse_invalidIP(t *testing.T) {
	_, err := client.GetDNSReverse([]string{"74.125.227", "63.11", "2747393"})

	assert.NotNil(t, err)
	_, ok := err.(*net.ParseError)
	assert.True(t, ok)
}
