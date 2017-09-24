package shodan

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_DeleteAlert(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	id := "ZZ4TDUUORVE1DIIP"
	path := fmt.Sprintf(alertDeletePath, id)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		fmt.Fprint(w, `{}`)
	})

	result, err := client.DeleteAlert(id)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestClient_GetAlert(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	id := "ZZ4TDUUORVE1DIIP"
	path := fmt.Sprintf(alertInfoPath, id)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "alert/alert"))
	})

	alert, err := client.GetAlert(id)
	alertExpected := &Alert{
		ID:         "ZZ4TDUUORVE1DIIP",
		Name:       "Test alert",
		Created:    "2017-09-24T18:30:43.592000",
		Expires:    0,
		Expired:    false,
		Expiration: "",
		Filters: &AlertFilters{
			IP: []string{"198.20.22.0/24"},
		},
		Size: 256,
	}

	assert.Nil(t, err)
	assert.Equal(t, alertExpected, alert)
}

func TestClient_GetAlerts(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(alertsInfoListPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "alert/alerts"))
	})

	alerts, err := client.GetAlerts()
	alertsExpected := []*Alert{
		{
			ID:         "ZZ4TDUUORVE1DIIP",
			Expired:    true,
			Name:       "Test alert",
			Created:    "2017-09-24T18:30:43.592000",
			Expires:    0,
			Expiration: "",
			Filters: &AlertFilters{
				IP: []string{"198.20.22.0/24"},
			},
			Size: 256,
		},
		{
			ID:         "IU0CJDXNNEXBOPK3",
			Name:       "Test alert 2",
			Expired:    false,
			Created:    "2017-09-24T20:08:51.815000",
			Expires:    100,
			Expiration: "2017-09-24T20:10:31.815000",
			Filters: &AlertFilters{
				IP: []string{"198.20.88.0/24"},
			},
			Size: 256,
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, alertsExpected, alerts)
}

func TestClient_CreateAlert(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(alertCreatePath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Write(getStub(t, "alert/create_alert"))
	})

	alert, err := client.CreateAlert("Test alert API", []string{"198.20.88.0/24"}, 0)
	alertExpected := &Alert{
		ID:         "JZT8NVWEZWCY79OO",
		Name:       "Test alert API",
		Created:    "2017-09-24T23:08:43.434646",
		Expires:    0,
		Expired:    false,
		Expiration: "",
		Filters: &AlertFilters{
			IP: []string{"198.20.88.0/24"},
		},
		Size: 256,
	}

	assert.Nil(t, err)
	assert.Equal(t, alertExpected, alert)
}
