package cmd

import (
	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

// albumsCmd represents the albums command
var albumsCmd = &cobra.Command{
	Use:   "albums",
	Short: "get all saved albums",
	Long:  `get all saved albums`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.UserAlbums(ctx, client, 1)
	},
}

func init() {
	rootCmd.AddCommand(albumsCmd)
}
