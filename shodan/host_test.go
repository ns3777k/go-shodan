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
	_, err := client.GetHostsForQuery(options)

	assert.Nil(t, err)
}

func TestHostVersion_UnmarshalJSON(t *testing.T) {
	payload := []byte(`{"vstr": "1.0", "vnum": 47}`)
	h := struct {
		VersionStr HostVersion `json:"vstr"`
		VersionNum HostVersion `json:"vnum"`
	}{}

	assert.Nil(t, json.Unmarshal(payload, &h))
	assert.Equal(t, "1.0", h.VersionStr.String())
	assert.Equal(t, "47", h.VersionNum.String())
}

func TestHostVersion_MarshalJSON(t *testing.T) {
	h := struct {
		VersionStr HostVersion `json:"v"`
	}{
		VersionStr: "25",
	}

	_, err := json.Marshal(h)
	assert.Nil(t, err)
}
