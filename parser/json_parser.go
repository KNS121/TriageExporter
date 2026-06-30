package parser

import (
	"bufio"
	"encoding/json"
	"os"
)

func ReadJSON(filePath string, handler func([]byte) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		err := handler(line)
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

func Parse(line []byte, eventType string) (map[string]interface{}, error) {
	var doc map[string]interface{}
	err := json.Unmarshal(line, &doc)
	if err != nil {
		return nil, err
	}

	doc["ModuleName"] = eventType

	//ApplyTransformations(doc, eventType)

	return doc, nil
}
