package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: "🙃0", expected: ""},
		{input: "🙂9", expected: "🙂🙂🙂🙂🙂🙂🙂🙂🙂"},
		{input: "aaф0b", expected: "aab"},

		{input: "১২৩", expected: "১২৩"},
		{input: "১2২৩0", expected: "১১২"},
		{input: "੩4", expected: "੩੩੩੩"},
		{input: `искус2тво`, expected: "искусство"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "1", expected: true},
		{input: "0", expected: true},
		{input: "9", expected: true},
		{input: "sd", expected: false},
		{input: " ", expected: false},
		{input: `\n`, expected: false},
		{input: `\\`, expected: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result := isDigit(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
