package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/common/dto"
	"net/http"
	"time"
)

// Client acts as an SDK so the Runner can communicate with the API
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

func (c *Client) UpdatePipelineRun(pipelineRunId string, updatePipelineRunDto dto.UpdatePipelineRunDto) (*pipelineruns.PipelineRun, error) {
	url := fmt.Sprintf("%s/runs/%s", c.BaseURL, pipelineRunId)

	payload, err := json.Marshal(updatePipelineRunDto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update pipeline run: %s", resp.Status)
	}

	var updatedPipelineRun pipelineruns.PipelineRun
	if err := json.NewDecoder(resp.Body).Decode(&updatedPipelineRun); err != nil {
		return nil, err
	}

	return &updatedPipelineRun, nil
}
