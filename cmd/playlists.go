package cmd

import (
	"gospt/internal/tui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playListsCmd)
}

var playListsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "Uses tui to show users playlists",
	Long:  `Opens tui showing all users playlists`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.DisplayPlaylists(ctx, client)
	},
}
