package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileOut, err := os.Create(toPath)
	if err != nil {
		return err
	}
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer file.Close()
	defer fileOut.Close()

	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}

		fileOut.Write(data[:n])
		// fmt.Print(string(data[:n]))
	}

	return nil
}
