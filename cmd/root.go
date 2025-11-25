package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/chzyer/readline"
	"github.com/misrab/clai/internal/ai"
	"github.com/spf13/cobra"
)

var (
	replMode bool
	aiModel  string
	useDummy bool

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
	rootCmd.Flags().StringVar(&aiModel, "model", "codellama:7b", "Ollama model to use")
	rootCmd.Flags().BoolVar(&useDummy, "dummy", false, "Use dummy AI (no Ollama required)")
	rootCmd.AddCommand(versionCmd)
}

// Execute wires stdout/stderr and runs the root command.
func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	return rootCmd.Execute()
}

// formatCommand returns the command with cyan coloring
func formatCommand(cmd string) string {
	return fmt.Sprintf("\033[36m%s\033[0m", cmd)
}

// handleSinglePrompt processes a single prompt in primary mode
func handleSinglePrompt(prompt string) error {
	command, err := generateCommand(prompt)
	if err != nil {
		return err
	}

	fmt.Printf("\nGenerated command:\n")
	fmt.Printf("  %s\n\n", formatCommand(command))

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

		command, err := generateCommand(prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Generated: %s\n", formatCommand(command))

		if err := promptAndExecute(command); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	return nil
}

// promptAndExecute asks for confirmation and executes the command
func promptAndExecute(command string) error {
	reader := bufio.NewReader(os.Stdin)

	for {
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
			rl, err := readline.NewEx(&readline.Config{
				Prompt:                 "Edit: ",
				InterruptPrompt:        "^C",
				HistoryLimit:           0,
				DisableAutoSaveHistory: true,
			})
			if err != nil {
				return err
			}
			// Prefill with current command
			rl.WriteStdin([]byte(command))
			edited, err := rl.Readline()
			rl.Close()
			if err != nil {
				if err == io.EOF || err == readline.ErrInterrupt {
					fmt.Println("Cancelled")
					return nil
				}
				return err
			}
			edited = strings.TrimSpace(edited)
			if edited != "" {
				command = edited
			}
			continue
		case "c", "copy":
			if err := clipboard.WriteAll(command); err != nil {
				return fmt.Errorf("copy failed: %w", err)
			}
			fmt.Println("Copied to clipboard")
			return nil
		default:
			fmt.Println("Invalid option, cancelled")
			return nil
		}
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

// generateCommand generates a shell command using AI or dummy mode
func generateCommand(prompt string) (string, error) {
	if useDummy {
		return generateDummyCommand(prompt), nil
	}

	client := ai.NewClient(aiModel)
	cmd, err := client.GenerateCommand(prompt)
	if err != nil {
		return "", err
	}
	return cmd, nil
}

// generateDummyCommand is a simple pattern-based command generator
func generateDummyCommand(prompt string) string {
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
