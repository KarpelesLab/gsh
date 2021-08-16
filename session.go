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
	scr := &Script{Filename: filename}
	return scr.run(s.getCtx(ctx), in)
}

func (s *Session) RunString(ctx context.Context, str string) error {
	reader := strings.NewReader(str)
	scr := &Script{Filename: "(inline)"}
	return scr.run(s.getCtx(ctx), reader)
}

func (s *Session) RunFile(ctx context.Context, filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	scr := &Script{Filename: filename}
	return scr.run(s.getCtx(ctx), fp)
}

type Job struct {
	os.Process
}
