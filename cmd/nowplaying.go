package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nowPlayingCmd)
}

var nowPlayingCmd = &cobra.Command{
	Use:     "nowplaying",
	Aliases: []string{"now"},
	Short:   "Shows song and artist of currently playing song",
	Long:    `Shows song and artist of currently playing song, useful for scripting`,
	Args:    cobra.MatchAll(cobra.RangeArgs(0, 1)),
	Run: func(cmd *cobra.Command, args []string) {
		commands.NowPlaying(ctx, args)
	},
}
