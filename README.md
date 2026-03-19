# Personal Terminal AI Helper (`help`)

A custom interactive Go-based terminal REPL (Read-Eval-Print Loop) that augments your daily workflow with an AI assistant. It allows you to execute standard shell commands and features an elegant Bubble Tea TUI for performing common AI-powered development tasks.

## Features

- **Interactive REPL**: Acts as a continuous terminal session.
- **AI Command Debugger**: Automatically reads the last failed command and its error output, dropping it into an AI context (`Ingenimax/agent-sdk-go`) to stream back the exact reason for the failure and the solution.
- **Smart Git Push**: Reads your `git diff` or staged changes, uses the AI to generate a clean, professional commit message, and automatically executes `git add .`, `git commit`, and `git push`.
- **Quick Build Menus**: A fast submenu for your most common build commands (`go run .`, `npm run build`, `docker-compose`, etc.).
- **Built-in Command History**: Saves your session history locally to `~/.help_history`.

## Installation

### Prerequisites
- Go 1.21+
- Ensure your `$GOPATH/bin` or `/usr/local/bin` is in your system's `$PATH`.

### Setup

1. **Clone or Navigate to the Repository:**
   ```bash
   cd personal-helper
   ```

2. **Configure your AI Credentials:**
   Copy the example configuration to create your active environment file:
   ```bash
   cp config/env.go.example config/env.go
   ```
   Open `config/env.go` and insert your actual `OPEN_ROUTER_KEY` and preferred `MODEL_ID`.

3. **Install the Binary:**
   Use the included Makefile to build and install the binary globally (this will also initialize your command history file):
   ```bash
   make install
   ```

## Usage

Simply type `help` in your terminal to enter the REPL:

```bash
$ help
help> 
```

- **Execute Commands**: Type any normal shell command (e.g., `ls -la`, `go run .`) and press `Enter`.
- **Open the AI Menu**: Press `Enter` on an empty `help>` prompt to launch the interactive Bubble Tea menu. From there, you can navigate using your arrow keys (or `j`/`k`) to select AI debugging, Git smart push, or build commands.
- **Exit**: Type `exit` or `quit`, or press `Ctrl+c`.

## Powered By
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI.
- [Ingenimax/agent-sdk-go](https://github.com/Ingenimax/agent-sdk-go) for AI integration and streaming.
