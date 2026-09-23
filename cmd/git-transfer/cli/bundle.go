package cli

import (
	"github.com/git-transfer/git-transfer/internal/capture"
	"github.com/spf13/cobra"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle [file.gtb]",
	Short: "Capture the current Git working state and create a bundle",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundlePath := args[0]
		return capture.Capture(bundlePath)
	},
}

func init() {
	rootCmd.AddCommand(bundleCmd)
}
