package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	aiModel         string
	useDummy        bool
	maxPromptLength int

	rootCmd = &cobra.Command{
		Use:   "clai",
		Short: "CLI for local AI",
		Long:  "clai - Use local AI for bash command generation and chat",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&aiModel, "model", "codellama:7b", "Ollama model to use")
	rootCmd.PersistentFlags().BoolVar(&useDummy, "dummy", false, "Use dummy AI (no Ollama required)")
	rootCmd.PersistentFlags().IntVar(&maxPromptLength, "max-length", 500, "Maximum prompt length in characters")
	rootCmd.AddCommand(versionCmd)
}

// validatePromptLength checks if prompt exceeds max length
func validatePromptLength(prompt string) error {
	if maxPromptLength > 0 && len(prompt) > maxPromptLength {
		return fmt.Errorf("prompt too long (%d chars). Max: %d chars (~1 paragraph)", len(prompt), maxPromptLength)
	}
	return nil
}

// Execute wires stdout/stderr and runs the root command.
func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	return rootCmd.Execute()
}
