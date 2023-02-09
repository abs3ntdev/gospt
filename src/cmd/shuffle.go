package cmd

import (
	"gitea.asdf.cafe/abs3nt/gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shuffleCmd)
}

var shuffleCmd = &cobra.Command{
	Use:   "shuffle",
	Short: "Toggles shuffle",
	Long:  `Enables shuffle if it is currently disabled or disables it if it is currently active`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Shuffle(ctx, client)
	},
}
