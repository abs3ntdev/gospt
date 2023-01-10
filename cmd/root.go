package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gospt/internal/auth"
	"gospt/internal/config"
	"gospt/internal/gctx"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

var (
	// Used for flags.
	ctx         *gctx.Context
	client      *spotify.Client
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "gospt",
		Short: "A spotify TUI and CLI to manage playback, browse library, and generate radios",
		Long:  `A spotify TUI and CLI to manage playback, borwse library, and generate radios written in go`,
	}
)

// Execute executes the root command.
func Execute(defCmd string) {
	if len(os.Args) == 1 {
		args := append([]string{defCmd}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	configDir, _ := os.UserConfigDir()
	cfgFile = filepath.Join(configDir, "gospt/client.yml")
	cobra.OnInitialize(initConfig)
	if !(len(os.Args) > 1) || os.Args[1] != "completion" {
		initConfig()
		var err error
		ctx = gctx.NewContext(context.Background())
		client, err = auth.GetClient(ctx)
		if err != nil {
			panic(err)
		}
		currentUser, err := client.CurrentUser(ctx)
		if err != nil {
			panic(err)
		}
		ctx.UserId = currentUser.ID
	}
}

func initConfig() {
	yamlDecoder := aconfigyaml.New()

	loader := aconfig.LoaderFor(&config.Values, aconfig.Config{
		AllowUnknownFields: true,
		AllowUnknownEnvs:   true,
		AllowUnknownFlags:  true,
		SkipFlags:          true,
		DontGenerateTags:   true,
		MergeFiles:         true,
		EnvPrefix:          "",
		FlagPrefix:         "",
		Files: []string{
			cfgFile,
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": yamlDecoder,
		},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
}
