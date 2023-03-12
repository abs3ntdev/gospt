package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.0.47"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints current verison",
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("Gospt: %s\n", Version)
}
