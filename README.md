# clai

CLI for local AI use - converts natural language prompts into shell commands and executes them with your approval.

## Requirements

- Go 1.21 or newer
- [Ollama](https://ollama.ai) (for AI features)

## Setup

### 1. Install Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.ai/install.sh | sh
```

### 2. Download a model

```bash
# Recommended: CodeLlama (fast, good for commands)
ollama pull codellama:7b

# Alternatives:
# ollama pull mistral
# ollama pull llama3.2
```

### 3. Start Ollama

```bash
# Start the Ollama service
ollama serve
```

Keep this running in a separate terminal, or run as a background service.

### 4. Install clai

```bash
# Build the binary
go build -o clai

# Or install globally
make install
```

Ensure `$HOME/go/bin` is in your PATH. Add to `~/.zshrc` if needed:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## Usage

### Primary Mode (Single-shot)

Ask clai to generate a command and approve it:

```bash
clai "copy all txt files to /tmp/backup"
```

Output:
```
Generated command:
  cp *.txt /tmp/backup/

Execute? [Y/n/e/c] 
```

Options:
- `Y` - Yes, execute the command
- `n` - No, cancel
- `e` - Edit the command inline before running
- `c` - Copy to clipboard

### REPL Mode (Interactive)

For multi-step workflows, use REPL mode:

```bash
clai --repl
```

Example session:
```
clai> copy all txt files to /tmp
Generated: cp *.txt /tmp/
Execute? [Y/n] y
✓ Executed

clai> show disk space
Generated: df -h
Execute? [Y/n] y
✓ Executed

clai> exit
Goodbye!
```

## Examples

```bash
# Basic usage (uses AI)
clai "find all txt files modified today"

# List files
clai "list files in current directory"

# Show disk usage
clai "show disk space"

# Use different model
clai --model mistral "compress my documents"

# Dummy mode (no AI, for testing without Ollama)
clai --dummy "list files"

# Interactive mode
clai --repl
```

## Flags

- `--repl` - Start in REPL (interactive) mode
- `--model <name>` - Specify Ollama model (default: `codellama:7b`)
- `--dummy` - Use pattern-based dummy mode (no Ollama required)

## Development

```bash
# Build
make build

# Run tests
make test

# Check version
clai version
```

## Notes

- AI features require Ollama to be running (`ollama serve`)
- Use `--dummy` flag to test without Ollama
- Models are cached locally by Ollama after first download
