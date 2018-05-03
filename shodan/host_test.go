package shodan

import (
	"net/http"
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetHostsForQuery_DifferentVersionFormats(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(hostSearchPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "host/version"))
	})

	options := &HostQueryOptions{Query: "argentina"}
	_, err := client.GetHostsForQuery(nil, options)

	assert.Nil(t, err)
}

func TestIntString_UnmarshalJSON(t *testing.T) {
	payload := []byte(`{"vstr": "1.0", "vnum": 47}`)
	h := struct {
		VersionStr IntString `json:"vstr"`
		VersionNum IntString `json:"vnum"`
	}{}

	assert.Nil(t, json.Unmarshal(payload, &h))
	assert.Equal(t, "1.0", h.VersionStr.String())
	assert.Equal(t, "47", h.VersionNum.String())
}

func TestIntString_MarshalJSON(t *testing.T) {
	h := struct {
		VersionStr IntString `json:"v"`
	}{
		VersionStr: "25",
	}

	_, err := json.Marshal(h)
	assert.Nil(t, err)
}

func TestHostIP_Unmarshal(t *testing.T) {
	testCases := []struct {
		expected    string
		jsonPayload []byte
	}{
		{"127.231.78.5", []byte(`
{"ip_str": "127.231.78.5", "ip": 3424324323, "data": [ {"ip_str": "127.231.78.5"} ]}`)},
		{"2600:1802:12::1", []byte(`
{"ip_str": "2600:1802:12::1", "ip": "2600:1802:12::1", "data": [ {"ip_str": "2600:1802:12::1"} ]}`)},
	}

	for _, testCase := range testCases {
		var h Host

		assert.Nil(t, json.Unmarshal(testCase.jsonPayload, &h))
		assert.Equal(t, testCase.expected, h.IP.String())
		assert.Equal(t, testCase.expected, h.Data[0].IP.String())
	}
}
