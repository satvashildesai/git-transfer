package cli

import (
	"github.com/git-transfer/git-transfer/internal/restore"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [file.gtb]",
	Short: "Check bundle integrity and repository compatibility without modifying working tree",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundlePath := args[0]
		return restore.Verify(bundlePath)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
