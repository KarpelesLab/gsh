package gsh

import (
	"io"
)

type Script struct {
	Filename string
}

type Command struct {
	s *Script
}

func (s *Script) run(ctx *Context, r io.Reader) error {
	// TODO
	return nil
}
