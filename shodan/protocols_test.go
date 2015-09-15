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
	})

	protocols, err := client.GetProtocols()
	if err != nil {
		t.Errorf("Client executeRequest returned error %v", err)
	}

	if len(protocols) != 116 {
		t.Errorf("There should be 116 protocols in stub %v", err)
	}
}
