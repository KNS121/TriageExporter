package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BulkDoc struct {
	Index    string
	Document map[string]interface{}
}

func (c *Client) BulkIndex(docs []BulkDoc) error {
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, doc := range docs {

		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": doc.Index,
			},
		}
		actionJSON, err := json.Marshal(action)
		if err != nil {
			return err
		}
		buf.Write(actionJSON)
		buf.WriteByte('\n')

		docJSON, err := json.Marshal(doc.Document)
		if err != nil {
			return err
		}
		buf.Write(docJSON)
		buf.WriteByte('\n')
	}

	url := fmt.Sprintf("%s/_bulk", c.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bulk request failed (status %d): %s", resp.StatusCode, string(data))
	}

	var result struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Index struct {
				Status int    `json:"status"`
				Error  string `json:"error,omitempty"`
			} `json:"index"`
		} `json:"items"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err == nil && result.Errors {
		errorCount := 0
		for _, item := range result.Items {
			if item.Index.Status >= 300 {
				errorCount++
			}
		}
		return fmt.Errorf("bulk indexing had %d errors", errorCount)
	}

	return nil
}
