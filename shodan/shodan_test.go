package shodan

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testClientToken = "TEST_TOKEN"
	stubsDir        = "stubs"
)

func setUpTestServe() (*http.ServeMux, func(), *Client) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(nil, testClientToken)
	client.BaseURL = server.URL
	client.ExploitBaseURL = server.URL
	client.StreamBaseURL = server.URL
	client.GeoNetBaseURL = server.URL
	tearDownTestServe := func() {
		server.Close()
	}

	return mux, tearDownTestServe, client
}

func getStub(t *testing.T, stubName string) []byte {
	t.Helper()
	stubPath := fmt.Sprintf("%s/%s.json", stubsDir, stubName)
	content, err := ioutil.ReadFile(stubPath)
	if err != nil {
		t.Errorf("getStub error %v", err)
	}

	return content
}

func TestNewClient(t *testing.T) {
	client := NewClient(nil, testClientToken)
	assert.Equal(t, testClientToken, client.Token)
}

func TestNewClient_httpClient(t *testing.T) {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	httpClient := &http.Client{Transport: transport}
	client := NewClient(httpClient, testClientToken)
	assert.ObjectsAreEqual(httpClient, client.Client)
}
