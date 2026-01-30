package main

import "os"

func ReadContent(fileName string) string {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}

	return string(content)
}
