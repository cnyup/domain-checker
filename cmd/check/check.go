package check

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xmapst/logx"

	"github.com/busybox-org/cert-checker/internal/core/checker"
	"github.com/busybox-org/cert-checker/internal/resolvers"
)

func New() *cobra.Command {
	root := &cobra.Command{
		Use:           "check",
		Short:         "Check the expiration date of the domain name",
		Long:          "Check the expiration date of the domain name",
		SilenceUsage:  true,
		SilenceErrors: true,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(resolvers.GetExternalIP())
			paths, err := cmd.Flags().GetStringSlice("path")
			if err != nil {
				logx.Fatalln(err)
			}
			suffix := cmd.Flags().Lookup("suffix").Value.String()
			days, err := cmd.Flags().GetInt("days")
			if err != nil {
				logx.Fatalln(err)
			}
			check := checker.New(suffix)
			res, err := check.CheckCerts(paths...)
			if err != nil {
				logx.Fatalln(err)
			}
			for _, v := range res {
				if v.ExpiredDays < 0 {
					_, _ = fmt.Fprintf(os.Stderr, "Path: %s, Doname:%s, ExpiredDay: %d, Is the domain name still valid!!!\n",
						v.Path, v.DomainName, v.ExpiredDays)
					continue
				}
				if v.ExpiredDays < days {
					_, _ = fmt.Fprintf(os.Stdout, "Path: %s, Doname:%s, ExpiredDay: %d\n",
						v.Path, v.DomainName, v.ExpiredDays)
				}
			}
			return nil
		},
	}
	return root
}
