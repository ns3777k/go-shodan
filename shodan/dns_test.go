package shodan

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"net"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetDNSResolve(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	expectedHostnames := []string{"google.com", "bing.com", "idonotexist.local"}

	mux.HandleFunc(resolvePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		hostnames := r.URL.Query().Get("hostnames")
		assert.NotEmpty(t, hostnames)

		split := strings.Split(hostnames, ",")
		assert.Len(t, split, len(expectedHostnames))

		w.Write(getStub(t, "dns_resolve")) //nolint:errcheck
	})

	resolve, err := client.GetDNSResolve(context.TODO(), expectedHostnames)

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

	expectedIPs := []net.IP{
		net.ParseIP("74.125.227.244"), net.ParseIP("92.63.108.40")}

	mux.HandleFunc(reversePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		ips := r.URL.Query().Get("ips")
		assert.NotEmpty(t, ips)

		split := strings.Split(ips, ",")
		assert.Len(t, split, len(expectedIPs))

		w.Write(getStub(t, "dns_reverse")) //nolint:errcheck
	})

	reversed, err := client.GetDNSReverse(context.TODO(), expectedIPs)

	assert.Nil(t, err)
	assert.Len(t, reversed, len(expectedIPs))

	for _, ip := range expectedIPs {
		_, ok := reversed[ip.String()]
		assert.True(t, ok)
	}
}
