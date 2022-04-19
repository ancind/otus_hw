package unpackstring

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
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
	tests := []struct {
		input    string
		expected string
	}{
		{input: "3sdvdg", expected: "element is number"},
		{input: "23423", expected: "element is number"},
		{input: "aaa10b", expected: "element is number"},
		{input: `\\234t23`, expected: "element is number"},
		{input: `1\3`, expected: "element is number"},
		{input: `vvv45`, expected: "element is number"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			_, err := Unpack(tc.input)
			require.EqualErrorf(t, err, tc.expected, "error message %v", err)
		})
	}
}
