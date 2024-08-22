package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of domain-checker",
	Long:  "All software has versions. This is domain-checker's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("domain-checker Static Site Generator v0.1 -- HEAD")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
