/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Print link to currently playing song",
	Long:  `Print link to currently playing song`,
	Run: func(cmd *cobra.Command, args []string) {
		link, err := commands.Link(ctx, client)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Print(link)
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}