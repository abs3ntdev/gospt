package gctx

import (
	"context"
	"fmt"
	"os"

	"tuxpa.in/a/zlog"
)

type Context struct {
	zlog.Logger
	Debug zlog.Logger

	context.Context
}

func (c *Context) Err() error {
	return c.Context.Err()
}

func (c *Context) Println(args ...any) {
	c.Info().Msg(fmt.Sprint(args...))
}

func NewContext(ctx context.Context) *Context {
	out := &Context{
		Context: ctx,
		Logger:  zlog.New(os.Stderr),
		Debug:   zlog.New(os.Stderr),
	}
	return out
}
