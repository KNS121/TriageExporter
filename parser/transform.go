package parser

import (
	"encoding/json"
)

func Parse(line []byte, eventType string) (map[string]interface{}, error) {

	var doc map[string]interface{}

	if err := json.Unmarshal(line, &doc); err != nil {
		return nil, err
	}

	doc["EventType"] = eventType

	return doc, nil
}
