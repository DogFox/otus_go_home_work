package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdList []string, env Environment) (returnCode int) {
	envForCmd := make([]string, 0)
	for key, value := range env {
		if value.NeedRemove {
			if err := os.Unsetenv(key); err != nil {
				return 1
			}
			continue
		}
		envForCmd = append(envForCmd, key+"="+value.Value)
		if err := os.Setenv(key, value.Value); err != nil {
			return 1
		}
	}
	envForCmd = append(envForCmd, os.Environ()...)

	cmd := exec.Command(cmdList[0], cmdList[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envForCmd

	if err := cmd.Run(); err != nil {
		return 1
	}
	return 0
}
