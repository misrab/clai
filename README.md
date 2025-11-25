# clai

CLI for local AI use - converts natural language prompts into shell commands and executes them with your approval.

## Requirements

- Go 1.21 or newer

## Installation

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
- `e` - Edit command first (coming soon)
- `c` - Copy to clipboard (coming soon)

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
# List files
clai "list files"

# Show disk usage
clai "show disk space"

# Compress files
clai "compress txt files"

# Interactive mode
clai --repl
```

## Development

```bash
# Build
make build

# Run tests
make test

# Check version
clai version
```

## Note

Currently using dummy AI responses for demonstration. Real AI integration coming soon!
