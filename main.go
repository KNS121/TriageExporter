package main

import (
	"TriageExporter/elasticsearch"
	"TriageExporter/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

const BULK_SIZE = 1000
const WORKERS = 4

func detectEventType(filePath string) string {
	name := filepath.Base(filePath)
	name = strings.ToLower(name)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	return name
}

func processFile(client *elasticsearch.Client, hostname, file string, wg *sync.WaitGroup,
	totalProcessed, totalIndexed, totalErrors, totalParseErrors *int64) {

	defer wg.Done()

	eventType := detectEventType(file)
	fmt.Printf("Processing: %s\n", file)

	var bulkDocs []elasticsearch.BulkDoc
	var mu sync.Mutex

	err := parser.ReadJSON(file, func(line []byte) error {
		atomic.AddInt64(totalProcessed, 1)

		doc, err := parser.Parse(line, eventType)
		if err != nil {
			atomic.AddInt64(totalParseErrors, 1)
			return nil
		}

		doc["SourceFile"] = filepath.Base(file)
		doc["Hostname"] = hostname

		index := "triage-" + eventType

		mu.Lock()
		bulkDocs = append(bulkDocs, elasticsearch.BulkDoc{
			Index:    index,
			Document: doc,
		})

		shouldFlush := len(bulkDocs) >= BULK_SIZE
		var docsToFlush []elasticsearch.BulkDoc
		if shouldFlush {
			docsToFlush = make([]elasticsearch.BulkDoc, len(bulkDocs))
			copy(docsToFlush, bulkDocs)
			bulkDocs = bulkDocs[:0]
		}
		mu.Unlock()

		if shouldFlush {
			err := client.BulkIndex(docsToFlush)
			if err != nil {
				atomic.AddInt64(totalErrors, int64(len(docsToFlush)))
				fmt.Printf("Bulk error: %v\n", err)
			} else {
				atomic.AddInt64(totalIndexed, int64(len(docsToFlush)))
			}
		}

		return nil
	})

	mu.Lock()
	if len(bulkDocs) > 0 {
		docsToFlush := bulkDocs
		bulkDocs = nil
		mu.Unlock()

		err := client.BulkIndex(docsToFlush)
		if err != nil {
			atomic.AddInt64(totalErrors, int64(len(docsToFlush)))
			fmt.Printf("Bulk error (final): %v\n", err)
		} else {
			atomic.AddInt64(totalIndexed, int64(len(docsToFlush)))
		}
	} else {
		mu.Unlock()
	}

	if err != nil {
		fmt.Printf("File processing error: %v\n", err)
	}
}

func main() {
	//if len(os.Args) < 2 {
	//	fmt.Println("Usage: TriageExporter <file|folder>")
	//	return
	//}

	//client := elasticsearch.New("http://localhost:9200")

	//fmt.Println("Loading templates...")
	//err := client.LoadTemplates("elasticsearch/templates")
	//if err != nil {
	//	fmt.Printf("Warning: failed to load templates: %v\n", err)
	//} else {
	//	fmt.Println("Templates loaded successfully")
	//}

	//path := os.Args[1]

	cfg := ReadConfig()

	client := elasticsearch.New(cfg.ELKURL)

	hostname := cfg.Hostname
	path := cfg.Path

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

	var totalProcessed int64
	var totalIndexed int64
	var totalErrors int64
	var totalParseErrors int64

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, WORKERS)

	for _, file := range files {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(f string) {
			defer func() { <-semaphore }()
			processFile(client, hostname, f, &wg,
				&totalProcessed, &totalIndexed, &totalErrors, &totalParseErrors)
		}(file)
	}

	wg.Wait()

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("TOTAL: processed=%d indexed=%d parse_errors=%d index_errors=%d\n",
		totalProcessed, totalIndexed, totalParseErrors, totalErrors)
	fmt.Println("Finished...")
}
