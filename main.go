package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gospt/authmanager"
	"gospt/config"
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
	ctx := context.Background()
	client, err := authmanager.GetClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	err = runner.Run(client, os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}
