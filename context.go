package gsh

import "context"

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
