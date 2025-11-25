# clai

`clai` is a tiny but production-minded "hello world" style CLI written in Go.
It demonstrates common best practices such as:

- modular project layout (`cmd/`, `internal/`)
- dependency management via Go modules
- Cobra-based command/flag parsing
- granular business logic that is easy to test
- version metadata that can be overridden at build time

## Requirements

- Go 1.21 or newer

## Getting Started

```bash
# run tests
make test

# build and run
make run ARGS="--name Ada"
make run ARGS="version"
```

## Installation

Install `clai` locally to use it anywhere:

```bash
make install
```

Ensure `$HOME/go/bin` is in your PATH. Add to `~/.zshrc` if needed:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## Usage

```bash
# default greeting
make run

# custom name
make run ARGS="--name Ada"

# version command
make run ARGS="version"

# help
make run ARGS="--help"
```

**Note:** The `ARGS` variable is needed because Make doesn't support passing flags directly (e.g., `make run --name Ada` won't work as Make would interpret `--name` as its own flag).
