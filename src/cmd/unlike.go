package cmd

import (
	"gitea.asdf.cafe/abs3nt/gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unlikeCmd)
}

var unlikeCmd = &cobra.Command{
	Use:     "unlike",
	Aliases: []string{"u"},
	Short:   "unlikes song",
	Long:    `unlikes song`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Unlike(ctx, client)
	},
}
