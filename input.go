package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const DefaultELKURL = "http://localhost:9200"

type Config struct {
	Hostname string
	Path     string
	ELKURL   string
}

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)

	text, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	return strings.TrimSpace(text)
}

func readHostname(reader *bufio.Reader) string {
	for {
		hostname := readLine(reader, "Enter hostname: ")

		if hostname == "" {
			fmt.Println("Hostname cannot be empty.\n")
			continue
		}

		return hostname
	}
}

func readPath(reader *bufio.Reader) string {
	for {
		path := readLine(reader, "Enter path to file/folder: ")

		if path == "" {
			fmt.Println("Path cannot be empty.\n")
			continue
		}

		if _, err := os.Stat(path); err != nil {
			fmt.Println("Path does not exist.\n")
			continue
		}

		return path
	}
}

func readELKURL(reader *bufio.Reader) string {
	url := readLine(reader, fmt.Sprintf("Enter ELK URL [%s]: ", DefaultELKURL))

	if url == "" {
		fmt.Printf("Using default ELK URL: %s\n\n", DefaultELKURL)
		return DefaultELKURL
	}

	return url
}

func ReadConfig() Config {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("==========================================")
	fmt.Println("           Triage Exporter")
	fmt.Println("==========================================")

	return Config{
		Hostname: readHostname(reader),
		Path:     readPath(reader),
		ELKURL:   readELKURL(reader),
	}
}
