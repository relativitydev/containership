package utils

import "strings"

// Before returns the substring before the first occurrence of a character
func Before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}

	return value[0:pos]
}

// After returns the substring after the first occurrence of a character
func After(value string, a string) string {
	// Get substring after a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}

	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}

	return value[adjustedPos:]
}
