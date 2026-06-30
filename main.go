package main

import (
	"TriageExporter/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: TriageExporter <file|folder>")
		return
	}

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

	for _, f := range files {
		fmt.Println("processing:", f)

		err := parser.ReadJSON(f, func(line []byte) error {
			//fmt.Println(string(line))
			return nil
		})

		if err != nil {
			fmt.Println("error reading file:", f, err)
		}
	}

	fmt.Println("\nREADY FOR NEXT STEP")
}
