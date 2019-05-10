package shodan

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Scan(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	expectedIPs := []string{"82.98.86.174", "93.184.216.34"}

	mux.HandleFunc(scanPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ips := r.FormValue("ips")
		assert.NotEmpty(t, ips)

		splited := strings.Split(ips, ",")
		assert.Len(t, splited, len(expectedIPs))

		for _, ip := range splited {
			assert.NotNil(t, net.ParseIP(ip))
		}

		w.Write(getStub(t, "scan")) //nolint:errcheck
	})

	scanStatus, err := client.Scan(context.TODO(), expectedIPs)
	scanStatusExpected := &CrawlScanStatus{
		ID:          "BOMA59VSGWX8QJR9",
		Count:       2,
		CreditsLeft: 183,
	}

	assert.Nil(t, err)
	assert.IsType(t, scanStatusExpected, scanStatus)
	assert.EqualValues(t, scanStatusExpected, scanStatus)
}

func TestClient_ScanInternet(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(scanInternetPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		port := r.FormValue("port")
		protocol := r.FormValue("protocol")

		assert.NotEmpty(t, port)
		assert.NotEmpty(t, protocol)

		_, err := strconv.Atoi(port)
		assert.Nil(t, err)

		fmt.Fprint(w, `{"id": "COMAD88STBX8QNN1"}`)
	})

	scanInternetStatusID, err := client.ScanInternet(context.TODO(), 22, "ssh")

	assert.Nil(t, err)
	assert.Equal(t, "COMAD88STBX8QNN1", scanInternetStatusID)
}

func TestClient_GetScanStatus(t *testing.T) {
	path := fmt.Sprintf(scanStatusPath, "BOMA59VSGWX8QJR9")
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "scan_status")) //nolint:errcheck
	})

	scanStatus, err := client.GetScanStatus(context.TODO(), "BOMA59VSGWX8QJR9")
	assert.Nil(t, err)

	scanStatusExpected := &ScanStatus{
		ID:     "BOMA59VSGWX8QJR9",
		Count:  2,
		Status: ScanStatusProcessing,
	}

	assert.IsType(t, scanStatusExpected, scanStatus)
	assert.EqualValues(t, scanStatusExpected, scanStatus)
}
