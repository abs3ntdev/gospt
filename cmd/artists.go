package cmd

import (
	"fmt"

	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

// artistsCmd represents the artists command
var artistsCmd = &cobra.Command{
	Use:   "artists",
	Short: "return all users artists",
	Long:  `return all users artists`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := commands.UserArtists(ctx, client, 1)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(artistsCmd)
}
