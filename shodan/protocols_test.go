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

	_, err := client.GetProtocols()
	if err != nil {
		t.Errorf("Client executeRequest returned error %v", err)
	}
}
