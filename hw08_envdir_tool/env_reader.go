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
		if strings.Contains(fileName.Name(), "=") { // скипнули файлы с =
			continue
		}

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
		needToRemove := false // будем удалять переменные с пустым значением
		if len(value) == 0 {
			needToRemove = true
		}
		str := prepareString(value)
		envs[fileName.Name()] = EnvValue{Value: str, NeedRemove: needToRemove}
	}

	return envs, nil
}

func prepareString(str string) string {
	strTrimmed := strings.TrimRight(str, " \t\r\n")
	finishStr := strings.ReplaceAll(strTrimmed, "\x00", "\n")
	return finishStr
}
