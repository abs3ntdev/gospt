package gctx

import (
	"context"
	"log"
	"os"
)

type Context struct {
	context.Context
	*log.Logger
	Debug  *log.Logger
	UserId string
}

func NewContext(ctx context.Context) *Context {
	out := &Context{
		ctx,
		log.New(os.Stdout, "LOG:", 0),
		log.New(os.Stdout, "DEBUG:", 1),
		"",
	}
	return out
}
