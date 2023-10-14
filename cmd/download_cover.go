package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCoverCmd)
}

var downloadCoverCmd = &cobra.Command{
	Use:     "download_cover",
	Aliases: []string{"now"},
	Short:   "Returns url for currently playing song art",
	Long:    `Returns url for currently playing song art`,
	Args:    cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		commands.DownloadCover(ctx, args)
	},
}
