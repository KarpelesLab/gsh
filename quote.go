package gsh

import (
	"bytes"
	"strconv"
	"strings"
)

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

// handleEscapes will process escape sequences in string, used by the echo builtin and the $'' escape
func handleEscapes(arg string) (string, bool) {
	buf := &bytes.Buffer{}

	for {
		p := strings.IndexByte(arg, '\\')
		if p == -1 {
			p = len(arg)
		}
		if p > 0 {
			buf.Write([]byte(arg[:p]))
			arg = arg[p:]
		}
		if len(arg) == 0 {
			return buf.String(), false
		}
		if arg[0] == '\\' {
			if len(arg) == 1 {
				buf.Write([]byte{'\\'})
				return buf.String(), false
			}
			v := arg[1]
			arg = arg[2:]
			fail := false
			switch v {
			case '\\':
				v = '\\'
			case 'a':
				v = '\a'
			case 'b':
				v = '\b'
			case 'c':
				return buf.String(), true // produce no further output
			case 'e':
				v = '\033'
			case 'f':
				v = '\f'
			case 'n':
				v = '\n'
			case 'r':
				v = '\r'
			case 't':
				v = '\t'
			case 'v':
				v = '\v'
			case '0':
				// read 3 more chars from arg
				if len(arg) >= 3 {
					oct, err := strconv.ParseInt(arg[:3], 8, 8)
					if err != nil {
						fail = true
						break
					}
					v = byte(oct)
				} else {
					fail = true
				}
			case 'x':
				if len(arg) >= 2 {
					ord, err := strconv.ParseInt(arg[:2], 16, 8)
					if err != nil {
						fail = true
						break
					}
					v = byte(ord)
				} else {
					fail = true
				}
			default:
				fail = true
			}
			if fail {
				buf.Write([]byte{'\\', v})
			} else {
				buf.Write([]byte{v})
			}

		}
	}
}
