package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	name string

	rootCmd = &cobra.Command{
		Use:   "clai",
		Short: "A tiny example CLI that follows Go best practices",
		Long:  "clai is a starter-friendly CLI that demonstrates structuring, flag handling, and versioning best practices in Go.",
		RunE: func(cmd *cobra.Command, args []string) error {
			message := greet(name)
			_, err := fmt.Fprintln(cmd.OutOrStdout(), message)
			return err
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "world", "Person or thing to greet")
	rootCmd.AddCommand(versionCmd)
}

// Execute wires stdout/stderr and runs the root command.
func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	return rootCmd.Execute()
}

func greet(target string) string {
	if target == "" {
		target = "world"
	}
	return fmt.Sprintf("Hello, %s!", target)
}
