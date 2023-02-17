package gctx

import (
	"context"
	"os"

	"tuxpa.in/a/zlog"
)

type Context struct {
	context.Context
	zlog.Logger
	Debug zlog.Logger
}

func NewContext(ctx context.Context) *Context {
	out := &Context{
		Context: ctx,
		Logger:  zlog.New(os.Stderr),
		Debug:   zlog.New(os.Stderr),
	}
	return out
}
