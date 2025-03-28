package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint: depguard
)

func TestReadDir(t *testing.T) {
	t.Run("try read dir", func(t *testing.T) {
		mapEnv, err := ReadDir("testdata/env")
		if err != nil {
			fmt.Println(fmt.Errorf("read env dir error: %w", err))
			return
		}
		require.Equal(t, mapEnv["EMPTY"].Value, "")
		require.Equal(t, mapEnv["UNSET"].NeedRemove, true)
		require.Equal(t, mapEnv["HELLO"].Value, "\"hello\"")
	})
}
