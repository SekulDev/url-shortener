package hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase10ToBase62(t *testing.T) {
	tests := []struct {
		base10 int64
		base62 string
	}{
		{0, "0"},
		{1, "1"},
		{61, "z"},
		{62, "10"},
		{124, "20"},
		{123456789, "8M0kX"},
	}

	for _, test := range tests {
		result := Base10ToBase62(test.base10)
		assert.Equal(t, test.base62, result, "Expected base62 conversion of %d to be %s, but got %s", test.base10, test.base62, result)
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"abc", "cba"},
	}

	for _, test := range tests {
		result := reverseString(test.input)
		assert.Equal(t, test.expected, result, "Expected reversal of %s to be %s, but got %s", test.input, test.expected, result)
	}
}
