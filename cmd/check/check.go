package check

import (
	"github.com/spf13/cobra"
)

var filePath string
var fileDir string

var suffix string
var days int

var accessToken string

var Check = &cobra.Command{
	Use:   "check",
	Short: "check cert",
	Long:  "check cert",
}

func init() {
	// 公用参数
	Check.PersistentFlags().StringVarP(&filePath, "path", "f", "", "cert file path")
	Check.PersistentFlags().StringVarP(&fileDir, "dir", "d", "", "cert file dir")
	Check.PersistentFlags().StringVarP(&suffix, "suffix", "", "", "file suffix match")
	Check.PersistentFlags().IntVarP(&days, "days", "", 15, "Alarm threshold days default 15")

	// dir 和 path 必须存在一个
	Check.MarkFlagsOneRequired("path", "dir")
	// path 和 regex 互斥  指定文件名时 不支持正则
	Check.MarkFlagsMutuallyExclusive("path", "suffix")
	//  path 和 dir 互斥  可以但没有必要指定两个
	Check.MarkFlagsMutuallyExclusive("path", "dir")

	Check.AddCommand(stdoutCmd)
	Check.AddCommand(dingTalkCmd)
}
