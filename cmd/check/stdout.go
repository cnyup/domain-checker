package check

import (
	"github.com/spf13/cobra"

	"github.com/cnyup/domain-checker/checker"
)

var stdoutCmd = &cobra.Command{
	Use:   "stdout",
	Short: "Print the result of stdout",
	Long:  "Print the result of stdout",
	RunE:  stdoutChecker,
}

func stdoutChecker(_ *cobra.Command, _ []string) error {
	c := checker.NewChecker()
	switch {
	case filePath != "":
		err := c.Check2StdoutByFile(filePath, days)
		if err != nil {
			return err
		}
	case fileDir != "":
		err := c.Check2StdoutByDir(fileDir, suffix, days)
		if err != nil {
			return err
		}
	}
	return nil
}
