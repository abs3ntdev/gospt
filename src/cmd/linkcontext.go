package cmd

import (
	"fmt"
	"os"

	"gitea.asdf.cafe/abs3nt/gospt/src/commands"

	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkContextCmd = &cobra.Command{
	Use:     "linkcontext",
	Aliases: []string{"lc"},
	Short:   "Get url to current context(album, playlist)",
	Long:    `Get url to current context(album, playlist)`,
	Run: func(cmd *cobra.Command, args []string) {
		link, err := commands.LinkContext(ctx, client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(link)
	},
}

func init() {
	rootCmd.AddCommand(linkContextCmd)
}
