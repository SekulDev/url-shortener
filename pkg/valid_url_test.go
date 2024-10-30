package pkg

import (
	"testing"
)

func TestIsValidUrl(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid cases
		{"https://youtube.com", true},
		{"http://10.1.18.20:443/test?a=b", true},
		{"https://example.com/path/to/page?query=1&foo=bar", true},
		{"http://localhost:8080", true},

		// Invalid cases
		{"youtube.com", false},       // Missing scheme
		{"ftp://example.com", false}, // Unsupported scheme
		{"http://", false},           // Missing host
		{"https://", false},          // Missing host
		{"", false},                  // Empty string
		{"http:/example.com", false}, // Malformed URL (missing one slash)
		{"://example.com", false},    // Missing scheme

		// Edge cases
		{"https://example.com#fragment", true}, // URL with fragment
		{"https://example.com/", true},         // URL with trailing slash
		{"https://123.123.123.123", true},      // IP address
	}

	for _, test := range tests {
		result := IsValidUrl(test.input)
		if result != test.expected {
			t.Errorf("IsValidUrl(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
