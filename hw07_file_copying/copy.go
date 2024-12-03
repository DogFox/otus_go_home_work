package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const bufferSize = 64

type copyInfo struct {
	fileSize   int64
	bufferSize int
	offset     int
	limit      int
}

func copyWorker(info copyInfo, file *os.File, fileOut *os.File) error {
	if info.fileSize < int64(info.offset) {
		return errors.New("offset bigger than file")
	}
	data := make([]byte, info.bufferSize)

	count := 0
	for {
		n, err := file.ReadAt(data, offset)
		offset += int64(n)
		if limit > 0 {
			if limit > int64(n) {
				limit -= int64(n)
			} else {
				fileOut.Write(data[:limit])
				fmt.Println("copied all 100%")
				break
			}
		}
		if err == io.EOF {
			fileOut.Write(data[:n])
			fmt.Println("copied all 100%")

			break
		}

		count++
		counts := int(info.fileSize) / info.bufferSize
		fmt.Println("copied [", strings.Repeat("#", count), strings.Repeat("-", counts-count), "]", count*100/(counts+1), "%")
		fileOut.Write(data[:n])
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileOut, err := os.Create(toPath)
	if err != nil {
		return err
	}
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	copyInfo := &copyInfo{
		fileSize:   fileInfo.Size(),
		bufferSize: bufferSize,
		limit:      int(limit),
		offset:     int(offset),
	}

	err = copyWorker(*copyInfo, file, fileOut)
	if err != nil {
		return err
	}

	defer file.Close()
	defer fileOut.Close()

	return nil
}
