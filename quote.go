package gsh

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
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
			v := rune(arg[1])
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
				ln := findLen(arg, "01234567", 3)
				log.Printf("ln = %d", ln)
				if ln == 0 {
					fail = true
					break
				}
				oct, err := strconv.ParseInt(arg[:ln], 8, 8)
				if err != nil {
					fail = true
					break
				}
				v = rune(oct)
				arg = arg[ln:]
			case 'x':
				ln := findLen(arg, "0123456789abcdefABCDEF", 2)
				if ln == 0 {
					fail = true
					break
				}
				ord, err := strconv.ParseInt(arg[:ln], 16, 8)
				if err != nil {
					fail = true
					break
				}
				v = rune(ord)
				arg = arg[ln:]
			case 'u':
				ln := findLen(arg, "0123456789abcdefABCDEF", 4)
				if ln == 0 {
					fail = true
					break
				}
				ord, err := strconv.ParseInt(arg[:ln], 16, 16)
				if err != nil {
					fail = true
					break
				}
				v = rune(ord)
				arg = arg[ln:]
			case 'U':
				ln := findLen(arg, "0123456789abcdefABCDEF", 8)
				if ln == 0 {
					fail = true
					break
				}
				ord, err := strconv.ParseInt(arg[:ln], 16, 32)
				if err != nil {
					fail = true
					break
				}
				v = rune(ord)
				arg = arg[ln:]
			default:
				fail = true
			}
			if fail {
				buf.WriteByte('\\')
				buf.WriteRune(v)
			} else {
				buf.WriteRune(v)
			}

		}
	}
}

func findLen(str, cutset string, max int) int {
	as, ok := makeASCIISet(cutset)
	if !ok {
		panic("string needs to be ascii only")
	}

	p := 0
	for {
		if p >= len(str) {
			return p
		}
		v := str[p]
		if !as.contains(v) {
			return p
		}
		p += 1
		if p >= max {
			return p
		}
	}
}

// Code below taken from: https://cs.opensource.google/go/go/+/refs/tags/go1.16.7:src/strings/strings.go;drc=refs%2Ftags%2Fgo1.16.7;l=804

// asciiSet is a 32-byte value, where each bit represents the presence of a
// given ASCII character in the set. The 128-bits of the lower 16 bytes,
// starting with the least-significant bit of the lowest word to the
// most-significant bit of the highest word, map to the full range of all
// 128 ASCII characters. The 128-bits of the upper 16 bytes will be zeroed,
// ensuring that any non-ASCII character will be reported as not in the set.
type asciiSet [8]uint32

// makeASCIISet creates a set of ASCII characters and reports whether all
// characters in chars are ASCII.
func makeASCIISet(chars string) (as asciiSet, ok bool) {
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		if c >= utf8.RuneSelf {
			return as, false
		}
		as[c>>5] |= 1 << uint(c&31)
	}
	return as, true
}

// contains reports whether c is inside the set.
func (as *asciiSet) contains(c byte) bool {
	return (as[c>>5] & (1 << uint(c&31))) != 0
}
