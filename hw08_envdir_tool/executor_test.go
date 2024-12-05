package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint: depguard
)

func TestRunCmd(t *testing.T) {
	t.Run("test run", func(t *testing.T) {
		mapEnv := make(map[string]EnvValue)

		mapEnv["EMPTY"] = EnvValue{
			Value:      "",
			NeedRemove: false,
		}

		mapEnv["UNSET"] = EnvValue{
			Value:      "NeedRemove",
			NeedRemove: true,
		}

		code := RunCmd([]string{"echo"}, mapEnv)

		_, okEmpty := os.LookupEnv("EMPTY")
		_, okUnset := os.LookupEnv("UNSET")

		require.Equal(t, 0, code)
		require.Equal(t, true, okEmpty)
		require.Equal(t, false, okUnset)
	})
}
