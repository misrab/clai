package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultOllamaURL = "http://localhost:11434"

// Client represents an Ollama API client
type Client struct {
	URL   string
	Model string
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

// NewClient creates a new Ollama client
func NewClient(model string) *Client {
	if model == "" {
		model = "codellama:7b"
	}
	return &Client{
		URL:   defaultOllamaURL,
		Model: model,
	}
}

// GenerateCommand converts a natural language prompt into a bash command
func (c *Client) GenerateCommand(prompt string) (string, error) {
	systemPrompt := fmt.Sprintf(`You are a bash command generator. Given a natural language request, output ONLY the bash command, nothing else.

Rules:
- Output ONLY the command, no explanations
- No markdown, no backticks, no formatting
- Single line command preferred
- Use common Unix/Linux commands

Request: %s

Command:`, prompt)

	reqBody := ollamaRequest{
		Model:  c.Model,
		Prompt: systemPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(c.URL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ollama not running? Install: https://ollama.ai")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	// Check for errors in response
	if ollamaResp.Error != "" {
		if strings.Contains(ollamaResp.Error, "not found") {
			return "", fmt.Errorf("model '%s' not found. Download it with:\n  ollama pull %s", c.Model, c.Model)
		}
		return "", fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	if ollamaResp.Response == "" {
		return "", fmt.Errorf("ollama returned empty response")
	}

	cmd := strings.TrimSpace(ollamaResp.Response)
	// Clean up common AI artifacts
	cmd = strings.Trim(cmd, "`")
	cmd = strings.TrimPrefix(cmd, "bash\n")
	cmd = strings.TrimPrefix(cmd, "sh\n")

	return cmd, nil
}
