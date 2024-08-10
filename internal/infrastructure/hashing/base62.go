package hashing

import (
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Base10ToBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	base := int64(len(base62Chars))
	var result strings.Builder

	for num > 0 {
		remainder := num % base
		result.WriteString(string(base62Chars[remainder]))
		num = num / base
	}

	// Reverse the string as the conversion gives us the result in reverse order
	return reverseString(result.String())
}

func reverseString(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
