package gsh

import (
	"io"
	"strings"
	"testing"
)

func TestTokens(t *testing.T) {
	sess := New()
	p := sess.newParser(strings.NewReader("echo $'This is a\\ntoken' and a few more\necho cmd2 ; echo cmd3"), "(test)")

	expect := []string{
		"echo|This is a\ntoken|and|a|few|more",
		"echo|cmd2",
		"echo|cmd3",
	}

	for {
		cmd, err := p.readCommand()

		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("failed, error in read: %s", err)
		}
		args, err := cmd.Argv(nil)
		if err != nil {
			t.Errorf("failed, error in argv: %s", err)
		}

		v := strings.Join(args, "|")
		if v != expect[0] {
			t.Errorf("failed, expected %q but got %q", expect[0], v)
		}
		expect = expect[1:]
	}
}
