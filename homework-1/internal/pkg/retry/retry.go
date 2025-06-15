package retry

import (
	"fmt"
	"net/http"
	"time"
)

type RetryClient struct {
	client     *http.Client
	maxRetries int
	delay      time.Duration
}

func New(maxRetries int, delay time.Duration) *RetryClient {
	return &RetryClient{
		client:     &http.Client{},
		maxRetries: maxRetries,
		delay:      delay,
	}
}

func (rc *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < rc.maxRetries; i++ {
		resp, err = rc.client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 420 && resp.StatusCode != 429 {
			return resp, nil
		}

		resp.Body.Close()
		time.Sleep(rc.delay)
	}

	return nil, fmt.Errorf("request failed after %d retries", rc.maxRetries)
}
