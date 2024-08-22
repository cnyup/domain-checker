package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/cnyup/domain-checker/cmd/check"
)

var rootCmd = &cobra.Command{
	Use:   "domain-checker",
	Short: "check domain cert",
	Long:  `check domain cert and send result to somewhere`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(check.Check)
}
