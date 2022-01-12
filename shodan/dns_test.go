package shodan

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetDomain(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	t1, _ := time.Parse(time.RFC3339Nano, "2019-09-30T08:49:55.868410+00:00")
	t2, _ := time.Parse(time.RFC3339Nano, "2019-09-30T08:49:55.867057+00:00")
	t3, _ := time.Parse(time.RFC3339Nano, "2019-10-06T17:01:04.237057+00:00")

	expected := &DomainDNSInfo{
		Domain:     "example.com",
		Tags:       []string{"ipv6", "spf"},
		Subdomains: []string{"www"},
		Data: []*SubdomainDNSInfo{
			{Subdomain: "", Type: "NS", Value: "b.iana-servers.net", LastSeen: t1},
			{Subdomain: "", Type: "NS", Value: "a.iana-servers.net", LastSeen: t2},
			{Subdomain: "www", Type: "AAAA", Value: "2606:2800:220:1:248:1893:25c8:1946", LastSeen: t3},
		},
	}

	path := fmt.Sprintf(dnsPath, "example.com")
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "dns_domain"))
	})

	actual, err := client.GetDomain(context.TODO(), "example.com")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClient_GetDNSResolve(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	expectedHostnames := []string{"google.com", "bing.com", "idonotexist.local"}

	mux.HandleFunc(resolvePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		hostnames := r.URL.Query().Get("hostnames")
		assert.NotEmpty(t, hostnames)

		split := strings.Split(hostnames, ",")
		assert.Len(t, split, len(expectedHostnames))

		w.Write(getStub(t, "dns_resolve"))
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
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	expectedIPs := []net.IP{net.ParseIP("74.125.227.244"), net.ParseIP("92.63.108.40")}

	mux.HandleFunc(reversePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		ips := r.URL.Query().Get("ips")
		assert.NotEmpty(t, ips)

		split := strings.Split(ips, ",")
		assert.Len(t, split, len(expectedIPs))

		w.Write(getStub(t, "dns_reverse"))
	})

	reversed, err := client.GetDNSReverse(context.TODO(), expectedIPs)

	assert.Nil(t, err)
	assert.Len(t, reversed, len(expectedIPs))

	for _, ip := range expectedIPs {
		_, ok := reversed[ip.String()]
		assert.True(t, ok)
	}
}
