package main

import "os"

func ModifyFile(fileName string, offset int, content string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Seek(int64(offset), 0)
	if err != nil {
		return
	}

	_, err = file.WriteString(content)
	if err != nil {
		return
	}
}
