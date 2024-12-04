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
				if exiterr, ok := err.(*exec.ExitError); ok {
					return exiterr.ExitCode()
				}
			}
			continue
		}
		envForCmd = append(envForCmd, key+"="+value.Value)
		if err := os.Setenv(key, value.Value); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				return exiterr.ExitCode()
			}
		}
	}
	envForCmd = append(envForCmd, os.Environ()...)

	cmd := exec.Command(cmdList[0], cmdList[1:]...) //nolint:gosec

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envForCmd

	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode()
		}
	}
	return 0
}
