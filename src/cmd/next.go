package cmd

import (
	"strconv"

	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nextCmd)
}

var nextCmd = &cobra.Command{
	Use:     "next",
	Aliases: []string{"n"},
	Args:    cobra.MatchAll(cobra.RangeArgs(0, 1)),
	Short:   "Skip to next song",
	Long:    `Skip to next song`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skipAmt := 1
		if len(args) >= 1 {
			var err error
			skipAmt, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		}
		return commands.Next(ctx, client, skipAmt)
	},
}
