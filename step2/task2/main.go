package main

import (
	"bufio"
	"io"
	"os"
)

func LineByNum(fileName string, lineNumber int) string {
	if lineNumber < 0 {
		return ""
	}

	file, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	defer file.Close()

	line := ""
	number := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if lineNumber == number {
			line = scanner.Text()
			break
		}

		number++
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return ""
	}

	return line
}
