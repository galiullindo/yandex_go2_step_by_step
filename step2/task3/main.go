package main

import (
	"io"
	"os"
)

func CopyFilePart(inFileName string, outFileName string, offset int) error {
	inFile, err := os.Open(inFileName)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = inFile.Seek(int64(offset), 0)
	if err != nil {
		return err
	}

	buf := make([]byte, 16)
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			_, err := outFile.Write(buf[:n])
			if err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
