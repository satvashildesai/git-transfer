package cli

import (
	"github.com/git-transfer/git-transfer/internal/restore"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [file.gtb]",
	Short: "Preview bundle contents",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundlePath := args[0]
		return restore.Inspect(bundlePath)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
