package gsh

import (
	"io"
	"log"
	"strings"
	"testing"
)

func TestTokens(t *testing.T) {
	sess := New()
	p := sess.newParser(strings.NewReader("echo 'This is a\ntoken' and a few more\necho cmd2 ; echo cmd3"), "(test)")

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
		log.Printf("cmd = %q", args)
	}
}
