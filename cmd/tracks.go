package cmd

import (
	"gospt/internal/tui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tracksCmd)
}

var tracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "Opens saved tracks",
	Long:  `Uses TUI to open a list of saved tracks`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.DisplayList(ctx, client)
	},
}
