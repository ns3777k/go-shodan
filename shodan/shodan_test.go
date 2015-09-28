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
	client = NewClient(testClientToken)
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
	client := NewClient(testClientToken)

	assert.Equal(t, testClientToken, client.Token)
}

func TestClient_buildURL(t *testing.T) {
	client := NewClient(testClientToken)
	testOptions := struct {
		Page    int  `url:"page"`
		ShowAll bool `url:"show_all"`
	}{
		100,
		true,
	}
	testCases := []struct {
		path     string
		params   interface{}
		expected string
	}{
		{
			"/testing/test/1",
			nil,
			baseURL + "/testing/test/1?key=" + testClientToken,
		},
		{
			"/testing/test/2",
			testOptions,
			baseURL + "/testing/test/2?key=" + testClientToken + "&page=100&show_all=true",
		},
	}

	for _, caseParams := range testCases {
		url, err := client.buildURL(baseURL, caseParams.path, caseParams.params)

		assert.Nil(t, err)
		assert.Equal(t, caseParams.expected, url)
	}
}

func TestClient_buildBaseURL(t *testing.T) {
	expected := client.BaseURL + "/test-base-url-building/?key=" + testClientToken
	actual, err := client.buildBaseURL("/test-base-url-building/", nil)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClient_buildExploitBaseURL(t *testing.T) {
	expected := client.ExploitBaseURL + "/test-exploit-url-building/?key=" + testClientToken
	actual, err := client.buildExploitBaseURL("/test-exploit-url-building/", nil)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClient_buildStreamBaseURL(t *testing.T) {
	expected := client.BaseURL + "/test-stream-url-building/?key=" + testClientToken
	actual, err := client.buildStreamBaseURL("/test-stream-url-building/", nil)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestClient_executeRequest_textUnauthorized(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	unauthorizedPath := "/http-error/401"

	errorText := "401 Unauthorized\n\n"
	errorText += "This server could not verify that you are authorized to access the document you requested.  " +
		"Either you supplied the wrong credentials (e.g., bad password), or your browser does not understand how to " +
		"supply the credentials required."

	mux.HandleFunc(unauthorizedPath, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, errorText, http.StatusUnauthorized)
	})

	url, err := client.buildBaseURL(unauthorizedPath, nil)
	assert.Nil(t, err)

	err = client.executeRequest("GET", url, nil, nil)
	assert.NotNil(t, err)
}

func TestClient_executeRequest_jsonNotFound(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	notFoundPath := "/http-error/404"

	mux.HandleFunc(notFoundPath, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error": "No information available for that IP."}`, http.StatusNotFound)
	})

	url, err := client.buildBaseURL(notFoundPath, nil)
	assert.Nil(t, err)

	err = client.executeRequest("GET", url, nil, nil)
	assert.NotNil(t, err)
}
