package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gospt/internal/auth"
	"gospt/internal/config"
	"gospt/internal/ctx"
	"gospt/internal/runner"
)

func init() {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config/gospt/")
	config.LoadConfig(configDir)
}

func main() {
	var err error
	log.New(os.Stdout, "LOG:", 0)
	ctx := ctx.NewContext(context.Background())
	client, err := auth.GetClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	err = runner.Run(ctx, client, os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}
