package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nextCmd)
}

var nextCmd = &cobra.Command{
	Use:     "next {amount}",
	Aliases: []string{"n", "skip"},
	Args:    cobra.MatchAll(cobra.RangeArgs(0, 1)),
	Short:   "Skip to next song or skip the specified number of tracks",
	Long:    `Skip to next song of skip the specified number of tracks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skipAmt := 1
		if len(args) >= 1 {
			var err error
			skipAmt, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		}
		return commands.Next(ctx, skipAmt)
	},
}
