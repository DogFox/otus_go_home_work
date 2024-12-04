package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdList []string, env Environment) (returnCode int) {

	fmt.Println(len(os.Args), os.Args)

	cmd := exec.Command(cmdList[0], cmdList[1:]...)
	cmd.Env = []string{} // не передаем какие либо енвы os.Ennviron

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return 1
	}
	return 0
}
