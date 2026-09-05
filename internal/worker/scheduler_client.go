package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type schedulerClient struct {
	http.Client
	URL string
}

func NewSchedulerClient(url string) *schedulerClient {
	return &schedulerClient{
		Client: http.Client{Timeout: 10 * time.Second},
		URL:    url,
	}
}

type errorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

func (c *schedulerClient) ReportSuccess(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.URL+"/jobs", nil)

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errResp); err != nil {
			return err
		}
		return fmt.Errorf("got error: %s -> %s", errResp.Error, errResp.Detail)
	}

	return nil
}
