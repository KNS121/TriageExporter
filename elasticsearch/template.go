package elasticsearch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (c *Client) LoadTemplates(templatesDir string) error {
	files, err := os.ReadDir(templatesDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			templatePath := filepath.Join(templatesDir, file.Name())
			templateData, err := os.ReadFile(templatePath)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", file.Name(), err)
			}

			templateName := strings.TrimSuffix(file.Name(), ".json")

			url := fmt.Sprintf("%s/_index_template/%s", c.BaseURL, templateName)

			req, err := http.NewRequest("PUT", url, bytes.NewReader(templateData))
			if err != nil {
				return fmt.Errorf("failed to create request for template %s: %w", templateName, err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := c.Client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to send template %s: %w", templateName, err)
			}

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				return fmt.Errorf("failed to load template %s (status %d): %s", templateName, resp.StatusCode, string(body))
			}
			resp.Body.Close()

			fmt.Printf("  Loaded template: %s\n", templateName)
		}
	}

	return nil
}
