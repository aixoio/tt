# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`tt` is a beautifully styled Git wrapper CLI built in Go using Cobra. It provides an intuitive, user-friendly interface for Git operations with interactive prompts, styled terminal output, and AI-powered features for commit messages and diff summaries.

## Key Technologies

- **CLI Framework**: Cobra (`github.com/spf13/cobra`) for command structure
- **Config Management**: Viper (`github.com/spf13/viper`) - config stored at `~/.tt/config.yaml`
- **Styling**: Lipgloss (`github.com/charmbracelet/lipgloss`) for terminal UI
- **Interactive Forms**: Huh (`github.com/charmbracelet/huh`) for prompts
- **Markdown Rendering**: Glamour (`github.com/charmbracelet/glamour`) for AI overview display
- **AI Integration**: OpenAI Go SDK v3 (`github.com/openai/openai-go/v3`) configured to use OpenRouter API

## Build and Run

```bash
# Build the binary
go build -o tt .

# Run directly
go run main.go <command>

# Run tests
go test ./cmd/...
```

## Architecture

### Command Structure
- Entry point: `main.go` → calls `cmd.Execute()`
- All commands live in `cmd/` directory as separate files
- Each command is a Cobra command registered with `rootCmd` in its `init()` function
- Root command defined in `cmd/root.go` with Viper config initialization

### Configuration System
- Config file: `~/.tt/config.yaml` (created automatically if missing)
- API key: Read from `TT_API_KEY` env var or `api_key` in config
- Settings:
  - `base_url`: AI API endpoint (default: `https://openrouter.ai/api/v1`)
  - `default_model`: Default AI model (default: `google/gemini-2.5-flash-lite`)
  - `diff_model`: Optional model override for `tt diff --ai` (falls back to `default_model`)
- Use `tt get <key>` and `tt set <key> <value>` to manage config

### Styling System
All styles defined in `styles/styles.go`:
- Color scheme: Deep blue primary, emerald success, red errors, violet highlights
- Pre-defined styles: `Header`, `Success`, `Error`, `Warning`, `Info`, `Highlight`, `Card`, etc.
- Icons: `SuccessIcon`, `ErrorIcon`, `WarningIcon`, `InfoIcon`
- Git-specific: `GitCommand`, `Branch`, `CommitHash`, `FilePath`
- Diff-specific: `DiffHeader`, `Add`, `Del`

### AI Features
Both `aic.go` and `diff.go` use OpenAI SDK with OpenRouter:
- **Commit message generation** (`tt aic`):
  - Analyzes git diff (staged or unstaged)
  - Includes project context (detects Go/Node.js/Java/Python/C++ via project files)
  - Includes changed file list
  - Interactive refinement loop: detailed, retry, summarize, feedback
  - Alias `tt a` auto-stages and auto-commits
- **Diff overview** (`tt diff --ai`):
  - Generates markdown summary of changes
  - Renders with Glamour for styled terminal output
  - Uses separate `diff_model` config if set

Client initialization pattern:
```go
client := openai.NewClient(
    option.WithBaseURL(baseURL),
    option.WithHeader("HTTP-Referer", "https://github.com/aixoio/tt"),
    option.WithHeader("X-Title", "tt"),
    option.WithAPIKey(apiKey),
)
```

### Utility Functions
- `runWithSpinner(title, action)`: Displays animated spinner during action
- `runWithSpinnerForMessage(title, action)`: Spinner that returns string result
- Both hide cursor during animation and restore after completion

## Command Implementation Patterns

### Standard Command Flow
1. Validate git repository with `git rev-parse --is-inside-work-tree`
2. Display styled header using `styles.Header.Render()`
3. Execute git commands via `exec.Command()`
4. Style output using appropriate `styles.*` constants
5. Show success/error with icons and styled messages

### Interactive Prompts
Use Huh forms with `huh.ThemeCharm()`:
```go
form := huh.NewForm(
    huh.NewGroup(
        huh.NewInput().Title("...").Value(&variable),
    ),
).WithTheme(huh.ThemeCharm())
form.Run()
```

### Adding New Commands
1. Create `cmd/<command>.go`
2. Define command with `var <name>Cmd = &cobra.Command{...}`
3. Implement `RunE: func(cmd *cobra.Command, args []string) error {...}`
4. Add to root in `init()`: `rootCmd.AddCommand(<name>Cmd)`
5. Use consistent styling from `styles/` package

## Important Notes

- All git operations should include error handling and user-friendly error messages
- Maintain consistent color scheme and styling patterns
- Use spinners for operations that may take time (especially AI calls)
- AI features require API key configuration via `tt set api_key <key>`
- The tool wraps git, so users should still have git installed
