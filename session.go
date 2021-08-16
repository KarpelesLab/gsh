package gsh

import (
	"context"
	"io"
	"os"
	"strings"
)

type Session struct {
	Location
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
