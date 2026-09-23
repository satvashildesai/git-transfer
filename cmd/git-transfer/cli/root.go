package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-transfer",
	Short: "A cross-platform CLI tool for safely packaging and transporting uncommitted Git working state",
	Long: `Git Transfer manages portable uncommitted working state.
It allows developers to transfer their current Git working state from one machine to another
without committing those changes on the source machine.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
