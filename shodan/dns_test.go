package shodan

import (
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetDNSResolve(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	hosts := []string{"google.com", "bing.com"}

	mux.HandleFunc(resolvePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		hostnames := r.URL.Query().Get("hostnames")
		assert.NotEmpty(t, hostnames)

		splited := strings.Split(hostnames, ",")
		assert.Len(t, splited, len(hosts))

		w.Write(getStub(t, "dns_resolve"))
	})

	resolve, err := client.GetDNSResolve(hosts)

	assert.Nil(t, err)
	assert.Len(t, resolve, len(hosts))

	for _, host := range hosts {
		ip, ok := resolve[host]
		assert.True(t, ok)
		assert.NotNil(t, net.ParseIP(ip))
	}
}
