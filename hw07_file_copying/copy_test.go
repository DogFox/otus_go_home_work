package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require" //nolint: depguard
)

func TestCopy(t *testing.T) {
	from := "./testdata/input.txt"
	to := "out.txt"
	var offset, limit int64 = 0, 0

	file, err := os.Open(from)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	out, err := os.Open(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	outInfo, err := out.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	t.Run("simple case", func(t *testing.T) {
		err = Copy(from, to, offset, limit)
		if err != nil {
			fmt.Println(err)
		}
		require.Equal(t, fileInfo.Size(), outInfo.Size())
	})
	cmd := exec.Command("rm", "-f", to)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
