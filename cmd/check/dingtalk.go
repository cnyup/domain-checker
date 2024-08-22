package check

import (
	"fmt"

	"github.com/spf13/cobra"

	"domainchecker/checker"
)

const CDingUrl = "https://oapi.dingtalk.com/robot/send?access_token"

var dingTalkCmd = &cobra.Command{
	Use:     "ding",
	Aliases: []string{"dingTalk", "dingtalk"},
	Short:   "Print the result to dingTalk",
	Long:    "Print the result to dingTalk",
	RunE:    dingTalkChecker,
}

func init() {
	dingTalkCmd.Flags().StringVarP(&accessToken, "token", "t", "", "dingTalk token (required)")
	_ = dingTalkCmd.MarkFlagRequired("token")

}

func dingTalkChecker(_ *cobra.Command, _ []string) error {
	var check = checker.NewChecker()
	url := fmt.Sprintf("%s=%s", CDingUrl, accessToken)
	check.Url = url

	switch {
	case filePath != "":
		err := check.Check2DingTalkByFile(filePath, days)
		if err != nil {
			return err
		}
	case fileDir != "":
		err := check.Check2DingTalkByDir(fileDir, suffix, days)
		if err != nil {
			return err
		}
	}
	return nil
}
