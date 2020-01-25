package command

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sawsh",
	Short: "Query and connect to ec2 instances",
	Args:  cobra.ExactArgs(1),
	Run:   connect,
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debugging")
	rootCmd.PersistentFlags().StringP("verbosity", "v", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	addConnectFlags(rootCmd)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := initLogs(os.Stderr, cmd.Flags().Lookup("verbosity").Value.String()); err != nil {
			return err
		}
		return nil
	}
}

func initLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
