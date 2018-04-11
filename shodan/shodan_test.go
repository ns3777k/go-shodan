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

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setUpTestServe() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient(nil, testClientToken)
	client.BaseURL = server.URL
	client.ExploitBaseURL = server.URL
	client.StreamBaseURL = server.URL
}

func getStub(t *testing.T, stubName string) []byte {
	stubPath := fmt.Sprintf("%s/%s.json", stubsDir, stubName)
	content, err := ioutil.ReadFile(stubPath)
	if err != nil {
		t.Errorf("getStub error %v", err)
	}

	return content
}

func tearDownTestServe() {
	server.Close()
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

//func TestClient_executeRequest_textUnauthorized(t *testing.T) {
//	setUpTestServe()
//	defer tearDownTestServe()
//
//	unauthorizedPath := "/http-error/401"
//
//	errorText := "401 Unauthorized\n\n"
//	errorText += "This server could not verify that you are authorized to access the document you requested.  " +
//		"Either you supplied the wrong credentials (e.g., bad password), or your browser does not understand how to " +
//		"supply the credentials required."
//
//	mux.HandleFunc(unauthorizedPath, func(w http.ResponseWriter, r *http.Request) {
//		http.Error(w, errorText, http.StatusUnauthorized)
//	})
//
//	url := client.buildBaseURL(unauthorizedPath, nil)
//	err := client.executeRequest("GET", url, nil, nil)
//
//	assert.NotNil(t, err)
//}
//
//func TestClient_executeRequest_jsonNotFound(t *testing.T) {
//	setUpTestServe()
//	defer tearDownTestServe()
//
//	notFoundPath := "/http-error/404"
//
//	mux.HandleFunc(notFoundPath, func(w http.ResponseWriter, r *http.Request) {
//		http.Error(w, `{"error": "No information available for that IP."}`, http.StatusNotFound)
//	})
//
//	url := client.buildBaseURL(notFoundPath, nil)
//	err := client.executeRequest("GET", url, nil, nil)
//
//	assert.NotNil(t, err)
//}
