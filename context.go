package gsh

import (
	"context"
	"io"
)

type Context struct {
	context.Context
	*Session
}

func getCtx(ctx context.Context) *Context {
	if v, ok := ctx.(*Context); ok {
		return v
	}
	return nil
}

func ctxOut(ctx *Context) io.Writer {
	return ctx.Session.Stdout
}
