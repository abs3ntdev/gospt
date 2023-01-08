package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gospt/authmanager"
	"gospt/config"
	"gospt/ctx"
	"gospt/runner"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
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
	client, err := authmanager.GetClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	err = runner.Run(ctx, client, os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}
