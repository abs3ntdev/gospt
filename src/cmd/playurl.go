package cmd

import (
	"gitea.asdf.cafe/abs3nt/gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playurlCmd)
}

var playurlCmd = &cobra.Command{
	Use:   "playurl",
	Short: "Plays song from provided url",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Long:  `Plays song from provided url`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.PlayUrl(ctx, client, args)
	},
}
