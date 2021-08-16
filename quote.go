package gsh

import "strings"

func Quote(s string) string {
	if len(s) == 0 {
		return "''"
	}
	// very simple quote process
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
