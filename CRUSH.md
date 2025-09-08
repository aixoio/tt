# CRUSH.md

## Build Commands
- Build the binary: `go build -o tt ./cmd`
- Install to $GOPATH/bin: `go install ./cmd`
- Clean: `go clean`

## Lint and Format Commands
- Format code: `gofmt -w .`
- Vet code: `go vet ./...`
- Static analysis: `go mod tidy && go mod verify`

## Code Style Guidelines
- **Formatting**: Always run `gofmt -w .` before committing. Use tabs for indentation (8 spaces equivalent).
- **Imports**: Group standard library imports first, then third-party, then local. Sort alphabetically within groups. Use `goimports -w .` for auto-management.
- **Naming Conventions**: 
  - Exported identifiers (functions, variables, types) start with uppercase (e.g., RootCmd).
  - Unexported start with lowercase (e.g., rootCmd).
  - Use camelCase for multi-word names (e.g., executeCommand).
  - Constants in UPPER_SNAKE_CASE.
- **Types and Interfaces**: Prefer interfaces over concrete types for dependencies. Use structs for data aggregation. Define types in the package where they're primarily used.
- **Error Handling**: Always check and propagate errors using `if err != nil { return err }`. Wrap errors with `fmt.Errorf("context: %w", err)` for context. In main, use `os.Exit(1)` for fatal errors.
- **Comments**: Use godoc-style comments for exported declarations (e.g., // Package cmd provides CLI commands). Keep comments concise and above the declaration.
- **Concurrency**: Use goroutines and channels for parallelism. Always sync with WaitGroup or select. Avoid shared state without mutexes.
- **Dependencies**: Use Cobra for CLI structure. Inject dependencies via structs. Avoid global variables.
- **General**: Follow effective Go guidelines. Keep functions short (&lt;50 lines). Use panic only for unrecoverable errors in libraries; prefer errors in apps.
- **Security**: Never log or expose sensitive data. Validate all inputs.

This file guides agentic coding agents in this Go CLI project using Cobra.
