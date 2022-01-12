package shodan

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetMyIP(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	testIP := "192.168.22.34"

	mux.HandleFunc(ipPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, strconv.Quote(testIP))
	})

	ip, err := client.GetMyIP(context.TODO())

	assert.Nil(t, err)
	assert.Equal(t, net.ParseIP(testIP), ip)
}

func TestClient_GetHTTPHeaders(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(headersPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "headers"))
	})

	headersExpected := map[string]string{
		"User-Agent":      "Go-http-client/1.1",
		"Host":            "api.shodan.io",
		"Accept-Encoding": "gzip",
	}
	headers, err := client.GetHTTPHeaders(context.TODO())

	assert.Nil(t, err)
	assert.Len(t, headers, len(headersExpected))
	assert.EqualValues(t, headersExpected, headers)
}
