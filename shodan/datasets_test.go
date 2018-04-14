package shodan

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestClient_GetDatasets(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	mux.HandleFunc(datasetsPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "datasets"))
	})

	datasets, err := client.GetDatasets(nil)
	assert.Nil(t, err)

	expectedDatasets := []*Dataset{{
		Name:        "raw-daily",
		Scope:       "daily",
		Description: "Data files containing all the information collected during a day"},
	}

	assert.Equal(t, expectedDatasets, datasets)
}

func TestClient_GetDatasetFiles(t *testing.T) {
	setUpTestServe()
	defer tearDownTestServe()

	path := fmt.Sprintf(datasetFilesPath, "raw-daily")
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.Write(getStub(t, "dataset_files"))
	})

	files, err := client.GetDatasetFiles(nil, "raw-daily")
	assert.Nil(t, err)

	expectedURL, _ := url.Parse("https://shodan.io/2017-12-29.json.gz")
	expectedFiles := []*DatasetFile{{
		Name:      "2017-12-29.json.gz",
		Size:      103750058939,
		Timestamp: time.Unix(1514669280, 0),
		URL:       expectedURL},
	}

	assert.Equal(t, expectedFiles, files)
}
