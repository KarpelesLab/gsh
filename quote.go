package gsh

import "strings"

func Quote(s string) string {
	if len(s) == 0 {
		return "''"
	}
	// very simple quote process
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

// isShellSafe returns true if character doesn't require quoting
func isShellSafe(c rune) bool {
	switch c {
	case '%', '+', '-', '.', '/', ':', '=', '@', '_':
		return true
	}

	if c >= '0' && c <= '9' {
		return true
	}
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	if c >= 0x80 {
		return true
	}
	return false
}
