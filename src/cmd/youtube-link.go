package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// youtubeLinkCmd represents the youtube-link command
var youtubeLinkCmd = &cobra.Command{
	Use:     "youtube-link",
	Aliases: []string{"yl"},
	Short:   "Print youtube link to currently playing song",
	Long:    `Print youtube link to currently playing song`,
	Run: func(cmd *cobra.Command, args []string) {
		link, err := commands.YoutubeLink(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(link)
	},
}

func init() {
	rootCmd.AddCommand(youtubeLinkCmd)
}
