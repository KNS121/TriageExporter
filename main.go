package main

import (
	"TriageExporter/elasticsearch"
	"TriageExporter/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func detectEventType(filePath string) string {
	name := filepath.Base(filePath)
	name = strings.ToLower(name)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	return name
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: TriageExporter <file|folder>")
		return
	}

	client := elasticsearch.New("http://localhost:9200")

	path := os.Args[1]

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	var files []string

	if info.IsDir() {
		fmt.Println("Input is folder:", path)
		filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".jsonl") {
				files = append(files, p)
			}
			return nil
		})
	} else {
		fmt.Println("Input is file:", path)
		files = append(files, path)
	}

	fmt.Printf("\nFound %d files:\n", len(files))
	for i, f := range files {
		info, err := os.Stat(f)
		if err == nil {
			fmt.Printf(" %d. %s (size: %d bytes)\n", i+1, f, info.Size())
		}
	}
	fmt.Println(strings.Repeat("-", 80))

	totalProcessed := 0
	totalIndexed := 0
	totalErrors := 0
	totalParseErrors := 0

	for _, file := range files {

		eventType := detectEventType(file)

		fmt.Println("Processing:", file)

		err := parser.ReadJSON(file, func(line []byte) error {

			doc, err := parser.Parse(line, eventType)
			if err != nil {
				return err
			}

			doc["SourceFile"] = filepath.Base(file)

			err = client.Index("triage", doc)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("TOTAL: processed=%d indexed=%d parse_errors=%d index_errors=%d\n",
		totalProcessed, totalIndexed, totalParseErrors, totalErrors)
	fmt.Println("READY FOR NEXT STEP")
}
