package main

import (
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/xmapst/logx"

	"github.com/busybox-org/cert-checker/cmd/check"
	"github.com/busybox-org/cert-checker/internal/core"
)

func main() {
	root := &cobra.Command{
		Use:           os.Args[0],
		Short:         "Check the expiration date of the domain name",
		Long:          "Check the expiration date of the domain name",
		SilenceUsage:  true,
		SilenceErrors: true,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logfile := cmd.Flags().Lookup("log_file").Value.String()
			if logfile != "" {
				logx.SetupConsoleLogger(logfile)
			}
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			logx.CloseLogger()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := filepath.Abs(os.Args[0])
			if err != nil {
				logx.Errorln(err)
				return err
			}
			svc, err := service.New(core.New(cmd.Flags()), &service.Config{
				Name:        name,
				DisplayName: name,
				Description: "Operating System Remote Executor Api",
				Executable:  name,
				Arguments:   os.Args[1:],
			})
			if err != nil {
				return err
			}
			err = svc.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	// check flags
	root.PersistentFlags().StringSliceP("path", "p", nil, "Directory or file paths to check (required)")
	_ = root.MarkPersistentFlagRequired("path")
	root.PersistentFlags().String("suffix", ".crt", "File suffix to check (Optional)")
	root.PersistentFlags().IntP("days", "d", 15, "Number of remaining days (Optional)")

	// alert flags
	root.Flags().StringP("alert_type", "t", "dingtalk", "Type of alert")
	root.Flags().String("alert_ak", "", "Access key for alerting (required)")
	_ = root.MarkFlagRequired("alert_ak")
	root.Flags().String("alert_sk", "", "Secret key for alerting (Optional)")
	// log flags
	root.Flags().StringP("log_file", "l", "", "Path to the log file (Optional)")
	// cron flags
	root.Flags().String("cron", "0 8  * * 1-5", "Cron expression for automatic execution (Optional)")
	// self update flags
	root.Flags().String("self_url", "https://oss.yfdou.com/tools/cert-checker", "URL for self-update (Optional)")

	root.AddCommand(
		check.New(),
	)
	if err := root.Execute(); err != nil {
		logx.Fatalln(err)
	}
}
