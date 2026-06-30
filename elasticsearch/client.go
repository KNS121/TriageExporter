package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Client  *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Index(index string, doc map[string]interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/_doc", c.BaseURL, index)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"elasticsearch returned %d: %s",
			resp.StatusCode,
			string(data),
		)
	}

	return nil
}
