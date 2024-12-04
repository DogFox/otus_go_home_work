package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadLine(r io.Reader, lineNum int) (line string, err error) {
	text, _ := bufio.NewReader(r).ReadString('\n')
	return text, nil

	// sc := bufio.NewScanner(r)
	// for sc.Scan() {
	// 	lastLine++
	// 	if lastLine == lineNum {
	// 		return sc.Text(), lastLine, sc.Err()
	// 	}
	// }
	// return line, lastLine, io.EOF
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	envs := make(map[string]EnvValue)
	for _, fileName := range files {

		file, err := os.Open(dir + "/" + fileName.Name())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		value, err := ReadLine(file, 1)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// fmt.Println(fileName, strings.Trim(value, "\n\r"))
		envs[file.Name()] = EnvValue{Value: strings.TrimRight(value, "\n\r"), NeedRemove: false}
	}

	return envs, nil
}
