package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	replMode bool

	rootCmd = &cobra.Command{
		Use:   "clai [prompt]",
		Short: "CLI for local AI - generate and execute shell commands",
		Long:  "clai converts natural language prompts into shell commands and executes them with your approval.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if replMode {
				return runREPL()
			}

			if len(args) == 0 {
				return fmt.Errorf("please provide a prompt (or use --repl for interactive mode)")
			}

			prompt := strings.Join(args, " ")
			return handleSinglePrompt(prompt)
		},
	}
)

func init() {
	rootCmd.Flags().BoolVar(&replMode, "repl", false, "Start in REPL (interactive) mode")
	rootCmd.AddCommand(versionCmd)
}

// Execute wires stdout/stderr and runs the root command.
func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	return rootCmd.Execute()
}

// handleSinglePrompt processes a single prompt in primary mode
func handleSinglePrompt(prompt string) error {
	command := generateCommand(prompt)

	fmt.Printf("\nGenerated command:\n")
	fmt.Printf("  \033[36m%s\033[0m\n\n", command)

	return promptAndExecute(command)
}

// runREPL starts the interactive REPL mode
func runREPL() error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("clai REPL mode - Type your requests (Ctrl+C to exit)")

	for {
		fmt.Print("\nclai> ")

		if !scanner.Scan() {
			break
		}

		prompt := strings.TrimSpace(scanner.Text())
		if prompt == "" {
			continue
		}

		if prompt == "exit" || prompt == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		command := generateCommand(prompt)
		fmt.Printf("Generated: \033[36m%s\033[0m\n", command)

		if err := promptAndExecute(command); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	return nil
}

// promptAndExecute asks for confirmation and executes the command
func promptAndExecute(command string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Execute? [Y/n/e/c] ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	response = strings.ToLower(strings.TrimSpace(response))

	switch response {
	case "", "y", "yes":
		return executeCommand(command)
	case "n", "no":
		fmt.Println("Cancelled")
		return nil
	case "e", "edit":
		fmt.Println("(Edit mode not yet implemented)")
		return nil
	case "c", "copy":
		fmt.Println("(Copy to clipboard not yet implemented)")
		fmt.Printf("Command: %s\n", command)
		return nil
	default:
		fmt.Println("Invalid option, cancelled")
		return nil
	}
}

// executeCommand runs the shell command
func executeCommand(command string) error {
	fmt.Println("Executing...")

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	fmt.Println("✓ Executed")
	return nil
}

// generateCommand is a dummy AI that generates shell commands
// TODO: Replace with actual AI integration
func generateCommand(prompt string) string {
	prompt = strings.ToLower(prompt)

	// Simple pattern matching for demo purposes
	switch {
	case strings.Contains(prompt, "copy") && strings.Contains(prompt, ".txt"):
		return "cp *.txt /tmp/backup/"
	case strings.Contains(prompt, "copy") && strings.Contains(prompt, "files"):
		return "cp -r ./files /tmp/backup/"
	case strings.Contains(prompt, "list"):
		return "ls -la"
	case strings.Contains(prompt, "disk"):
		return "df -h"
	case strings.Contains(prompt, "compress") || strings.Contains(prompt, "zip"):
		return "tar -czf backup.tar.gz *.txt"
	case strings.Contains(prompt, "delete") || strings.Contains(prompt, "remove"):
		return "rm -i unwanted_file.txt"
	default:
		return fmt.Sprintf("echo 'Dummy command for: %s'", prompt)
	}
}
