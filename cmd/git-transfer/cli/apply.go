package cli

import (
	"github.com/git-transfer/git-transfer/internal/restore"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply [file.gtb]",
	Short: "Validate the bundle against the current repository and restore its working state",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundlePath := args[0]
		return restore.Apply(bundlePath)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
