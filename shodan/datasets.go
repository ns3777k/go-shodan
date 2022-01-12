package shodan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const (
	datasetsPath     = "/shodan/data"
	datasetFilesPath = "/shodan/data/%s"
)

// Dataset represents short information about single dataset.
type Dataset struct {
	Name        string `json:"name"`
	Scope       string `json:"scope"`
	Description string `json:"description"`
}

// DatasetFile contains files to a dataset.
type DatasetFile struct {
	URL       *url.URL  `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
}

// UnmarshalJSON implements Unmarshaler interface for 2 reasons:
//
// 1. transform int timestamp into valid time.Time.
// 2. transform url string into *url.URL.
func (f *DatasetFile) UnmarshalJSON(data []byte) error {
	type Alias DatasetFile

	ff := &struct {
		*Alias
		Timestamp int64  `json:"timestamp"`
		URL       string `json:"url"`
	}{
		Alias: (*Alias)(f),
	}

	err := json.Unmarshal(data, &ff)
	if err != nil {
		return err
	}

	u, err := url.Parse(ff.URL)
	if err != nil {
		return err
	}

	f.Timestamp = time.Unix(ff.Timestamp/1000, 0)
	f.URL = u

	return nil
}

// MarshalJSON implements Marshaler interface to restore original json.
//
// 1. transform time.Time into int.
// 2. transform *url.URL into string.
func (f *DatasetFile) MarshalJSON() ([]byte, error) {
	type Alias DatasetFile

	return json.Marshal(&struct {
		*Alias
		Timestamp int64  `json:"timestamp"`
		URL       string `json:"url"`
	}{
		Alias:     (*Alias)(f),
		Timestamp: f.Timestamp.Unix() * 1000,
		URL:       f.URL.String(),
	})
}

// GetDatasets provides list of the datasets that are available for download.
func (c *Client) GetDatasets(ctx context.Context) ([]*Dataset, error) {
	datasets := make([]*Dataset, 0)
	req, err := c.NewRequest("GET", datasetsPath, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &datasets); err != nil {
		return nil, err
	}

	return datasets, nil
}

// GetDatasetFiles returns a list of files that are available for download from the provided dataset.
func (c *Client) GetDatasetFiles(ctx context.Context, name string) ([]*DatasetFile, error) {
	files := make([]*DatasetFile, 0)
	path := fmt.Sprintf(datasetFilesPath, name)

	req, err := c.NewRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Do(ctx, req, &files); err != nil {
		return nil, err
	}

	return files, nil
}
