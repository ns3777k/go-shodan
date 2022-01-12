package shodan

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetAlertTriggers(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(alertTriggersListPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "alert/triggers"))
	})

	triggers, err := client.GetAlertTriggers(context.TODO())
	triggersExpected := []*AlertTrigger{
		{
			Name:        "any",
			Rule:        "*",
			Description: "Match any service that is discovered",
		},
		{
			Name:        "industrial_control_system",
			Rule:        "tag:ics",
			Description: "Services associated with industrial control systems",
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, triggersExpected, triggers)
}

func TestClient_EnableAlertTrigger(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	id := "TestClient_EnableAlertTrigger"
	trigger := "malware"
	path := fmt.Sprintf(alertTriggerEnablePath, id, trigger)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		fmt.Fprint(w, `{"success": true}`)
	})

	options := &AlertTriggerIdent{AlertID: id, TriggerName: trigger}
	r, err := client.EnableAlertTrigger(context.TODO(), options)

	assert.Nil(t, err)
	assert.True(t, r)
}

func TestClient_DisableAlertTrigger(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	id := "TestClient_DisableAlertTrigger"
	trigger := "ssl_expired"
	path := fmt.Sprintf(alertTriggerEnablePath, id, trigger)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		fmt.Fprint(w, `{"success": true}`)
	})

	options := &AlertTriggerIdent{AlertID: id, TriggerName: trigger}
	r, err := client.DisableAlertTrigger(context.TODO(), options)

	assert.Nil(t, err)
	assert.True(t, r)
}

func TestClient_AddServiceToAlertTriggerWhitelist(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	id := "TestClient_AddServiceToAlertTriggerWhitelist"
	trigger := "iot"
	service := "127.0.0.1:80"
	path := fmt.Sprintf(alertTriggerWhitelistPath, id, trigger, service)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		fmt.Fprint(w, `{"success": true}`)
	})

	triggerOptions := &AlertTriggerIdent{AlertID: id, TriggerName: trigger}
	serviceOptions := &AlertTriggerServiceIdent{ServiceName: service, AlertTriggerIdent: triggerOptions}
	r, err := client.AddServiceToAlertTriggerWhitelist(context.TODO(), serviceOptions)

	assert.Nil(t, err)
	assert.True(t, r)
}

func TestClient_RemoveServiceFromAlertTriggerWhitelist(t *testing.T) {
	mux, tearDownTestServe, client := setUpTestServe()
	defer tearDownTestServe()

	id := "TestClient_RemoveServiceFromAlertTriggerWhitelist"
	trigger := "open_database"
	service := "127.0.0.1:80"
	path := fmt.Sprintf(alertTriggerWhitelistPath, id, trigger, service)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		fmt.Fprint(w, `{"success": true}`)
	})

	triggerOptions := &AlertTriggerIdent{AlertID: id, TriggerName: trigger}
	serviceOptions := &AlertTriggerServiceIdent{ServiceName: service, AlertTriggerIdent: triggerOptions}
	r, err := client.RemoveServiceFromAlertTriggerWhitelist(context.TODO(), serviceOptions)

	assert.Nil(t, err)
	assert.True(t, r)
}
