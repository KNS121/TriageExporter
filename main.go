package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: TriageExporter.exe <file|folder>")
		fmt.Println("Example: TriageExporter.exe Process.jsonl")
		return
	}

	path := os.Args[1]

	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	if info.IsDir() {
		fmt.Println("Input is folder:", path)
	} else {
		fmt.Println("Input is file:", path)
	}

	fmt.Println("READY")
}
