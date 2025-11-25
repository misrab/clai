package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/misrab/clai/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), version.Full())
	},
}
