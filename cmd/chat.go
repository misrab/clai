package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/misrab/clai/internal/ai"
	"github.com/spf13/cobra"
)

var (
	chatNoRepl bool

	chatCmd = &cobra.Command{
		Use:   "chat [prompt]",
		Short: "Chat with AI (always REPL mode)",
		Long:  "Have a conversation with the AI without command execution. Always starts in REPL mode. Provide an initial prompt to auto-submit it as the first message.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var initialPrompt string
			if len(args) > 0 {
				initialPrompt = strings.Join(args, " ")
			}

			// Check if --no-repl flag is set for single-shot mode
			if chatNoRepl {
				if initialPrompt == "" {
					return fmt.Errorf("please provide a prompt for single-shot mode")
				}
				return handleChatPrompt(initialPrompt)
			}

			// Always REPL mode (default)
			return runChatREPL(initialPrompt)
		},
	}
)

func init() {
	chatCmd.Flags().BoolVar(&chatNoRepl, "no-repl", false, "Single-shot mode instead of REPL")
	rootCmd.AddCommand(chatCmd)
}

// handleChatPrompt processes a single chat prompt (--no-repl mode)
func handleChatPrompt(prompt string) error {
	if useDummy {
		fmt.Printf("Dummy response to: %s\n", prompt)
		return nil
	}

	client := ai.NewClient(aiModel)
	response, err := client.Chat(prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}
	fmt.Println(response)
	return nil
}

// runChatREPL starts the interactive chat REPL mode
func runChatREPL(initialPrompt string) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\033[2mclai chat - Type your messages ('exit' to quit)\033[0m")

	if initialPrompt != "" {
		fmt.Printf("\033[1;34mYou:\033[0m %s\n", initialPrompt)
		if err := validatePromptLength(initialPrompt); err != nil {
			fmt.Printf("\033[31m%v\033[0m\n", err)
		} else if err := streamChatResponse(initialPrompt); err != nil {
			fmt.Printf("\033[31mError: %v\033[0m\n", err)
		}
	}

	for {
		fmt.Print("\n\033[1;34mYou:\033[0m ")
		if !scanner.Scan() {
			break
		}

		prompt := strings.TrimSpace(scanner.Text())
		if prompt == "" {
			continue
		}
		if prompt == "exit" || prompt == "quit" {
			fmt.Println("\033[2mGoodbye!\033[0m")
			break
		}

		if err := validatePromptLength(prompt); err != nil {
			fmt.Printf("\033[31m%v\033[0m\n", err)
			continue
		}
		if err := streamChatResponse(prompt); err != nil {
			fmt.Printf("\033[31mError: %v\033[0m\n", err)
			continue
		}
	}

	return nil
}

// streamChatResponse streams the AI response
func streamChatResponse(prompt string) error {
	if useDummy {
		fmt.Printf("\n\033[1;32mAI:\033[0m Dummy response to: %s\n", prompt)
		return nil
	}

	client := ai.NewClient(aiModel)
	fmt.Print("\n\033[1;32mAI:\033[0m ")

	err := client.ChatStream(prompt, func(chunk string) error {
		fmt.Print(chunk)
		return nil
	})

	fmt.Println()
	return err
}
