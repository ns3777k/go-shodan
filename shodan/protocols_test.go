package shodan

import (
	"testing"
	"net/http"
)

func TestClient_GetProtocols(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(protocolsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Write(getStub(t, "protocols"))
	});

	url, err := client.buildURL(server.URL, protocolsPath, nil)
	if err != nil {
		t.Errorf("buildURL returned error %v", err)
	}

	var protocols ProtocolCollection
	err = client.executeRequest("GET", url, &protocols)
	if err != nil {
		t.Errorf("Client executeRequest returned error %v", err)
	}
}
