package strings

import "strings"

// IsBlank Checks if the given string has text.
//
//	isBlank("")        = true
//	isBlank(" ")       = true
//	isBlank("batman")     = false
//	isBlank("  batman  ") = false
func IsBlank(s string) bool {
	if len(strings.TrimSpace(s)) == 0 {
		return true
	}

	return false
}

func IsNotBlank(s string) bool {
	return !IsBlank(s)
}
