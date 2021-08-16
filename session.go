package gsh

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"
)

type Session struct {
	Location
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func New() *Session {
	return &Session{
		Location: ProcLocation(),
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
	}
}

func (s *Session) getCtx(ctx context.Context) *Context {
	if ctx == nil {
		return &Context{Context: context.Background(), Session: s}
	}
	if v, ok := ctx.(*Context); ok {
		if v.Session == s {
			return v
		}
	}
	return &Context{Context: ctx, Session: s}
}

func (s *Session) newParser(r io.Reader, filename string) *parser {
	p := &parser{session: s, line: 1, col: 1, filename: filename}
	// if v is already a bufio.Reader, use it as is
	switch v := r.(type) {
	case *bufio.Reader:
		p.bio = v
	default:
		p.bio = bufio.NewReader(v)
	}
	return p
}

func (s *Session) Run(ctx context.Context, in io.Reader, filename string) error {
	return runScript(s.getCtx(ctx), in, filename)
}

func (s *Session) RunString(ctx context.Context, str string) error {
	reader := strings.NewReader(str)
	return runScript(s.getCtx(ctx), reader, "(inline)")
}

func (s *Session) RunFile(ctx context.Context, filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	return runScript(s.getCtx(ctx), fp, filename)
}

type Job struct {
	os.Process
}
