package main

import (
	"fmt"
	"os"
)

// go run . "./testdata/env" "/bin/bash" "./testdata/echo.sh" arg1=1 arg2=2.
func main() {
	dirPath := os.Args[1]  // забрал дирректорию с енвами
	cmdArgs := os.Args[2:] // остальное - команда с ее инпутом, не парсим
	envs, err := ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(envs)
	RunCmd(cmdArgs, envs)
}
