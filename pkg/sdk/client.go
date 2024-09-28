package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spyrosmoux/api/pkg/business/pipelineruns"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) UpdatePipelineRunStatus(pipelineRunId string, status pipelineruns.Status) (*pipelineruns.PipelineRun, error) {
	url := fmt.Sprintf("%s/runs/%s", c.BaseURL, pipelineRunId)

	payload, err := json.Marshal(status.String())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Make the HTTP request to the API
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update pipeline status: %s", resp.Status)
	}

	// Decode the response body into a User struct
	var updatedPipelineRun pipelineruns.PipelineRun
	if err := json.NewDecoder(resp.Body).Decode(&updatedPipelineRun); err != nil {
		return nil, err
	}

	return &updatedPipelineRun, nil
}
