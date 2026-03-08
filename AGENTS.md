# AGENTS.md - Console Text Editor Project

## Project Overview

Go project implementing a console-based text editor with SSH support.
Main function signature:
```go
func consoleEditText(defaultTextA string, optsA ...string) map[string]interface{}
```

## Build Commands

```bash
go build -o console_editor .
GOOS=linux GOARCH=amd64 go build -o console_editor_linux .
GOOS=windows GOARCH=amd64 go build -o console_editor.exe .
go run . -filePath=/path/to/file.txt
```

## Test Commands

```bash
go test ./...
go test -run TestFunctionName ./...
go test -v ./...
go test -cover ./...
go test -bench=. ./...
```

## Linting & Code Quality

```bash
golint ./...
go vet ./...
staticcheck ./...
go vet && golint ./... && staticcheck ./...
go fmt ./...
```

## Code Style Guidelines

### General Principles

- Use Go standard library whenever possible; avoid external dependencies
- Write clear, idiomatic Go code following official Go conventions
- Keep functions small and focused (single responsibility)

### Naming Conventions

- **Variables**: `camelCase` (e.g., `defaultText`, `filePath`)
- **Constants**: `PascalCase` (e.g., `MaxFileSize`)
- **Functions**: `PascalCase` (e.g., `consoleEditText`)
- **Types/Structs**: `PascalCase` (e.g., `EditorState`)
- **Packages**: `lowercase`, short (e.g., `editor`, `ssh`)
- **Interfaces**: `PascalCase` with `er` suffix (e.g., `Reader`)

### Imports

Group imports in three sections: standard library, external packages, internal packages.

### Formatting

- Use `go fmt` for automatic formatting
- Maximum line length: ~100 characters (soft limit)
- Use tabs for indentation
- Add blank line between top-level declarations

### Types & Type Safety

- Use explicit types; avoid `var x` without type
- Prefer specific types over `interface{}` when possible
- Use custom types for domain concepts

### Error Handling

- Always handle errors explicitly; no ignored errors (`_`)
- Return errors with context: `fmt.Errorf("context: %w", err)`
- Check errors early and return/fatal early
- Wrap errors at package boundaries

### Function Design

- Keep functions under 50 lines when possible
- Maximum 3-4 parameters; use structs for many parameters
- Return early, avoid deeply nested conditionals
- Document exported functions with Go doc comments

### Concurrency

- Use goroutines and channels for concurrent operations
- Always handle channel closure
- Use `sync.WaitGroup` for coordinating goroutines
- Pass context for cancellation

### Testing

- Test files: `*_test.go` in same package
- Use table-driven tests for multiple cases
- Name tests descriptively: `TestConsoleEditText_WithSSHConfig`

### UTF-8 & Chinese Support

- All strings must be UTF-8 encoded
- Use `unicode/utf8` package for rune-based operations
- Handle multi-byte characters correctly (not byte-oriented)

### SSH Features

- Support password and key-based authentication
- Use `golang.org/x/crypto/ssh`
- Handle connection timeouts gracefully

### Editor Features

- Keyboard shortcuts: Ctrl+S (save), Ctrl+K (save as), Ctrl+C/V (copy/paste)
- Ctrl+Z/Y (undo/redo), Ctrl+F (find), Ctrl+H (replace), Ctrl+G (goto line)
- Ctrl+X (exit), Ctrl+Q (force quit), Ctrl+W (toggle wrap)
- Find/replace supports regex
- Auto-wrap text by default
- Return map with keys: "text", "status", "error" (on error)

## Project Structure

```
.
├── main.go              # Entry point
├── editor/
│   ├── editor.go        # Core editor logic
│   └── editor_test.go
├── ssh/
│   ├── ssh.go           # SSH file operations
│   └── ssh_test.go
├── utils/
│   └── utils.go         # Helper functions
└── go.mod               # Go module file
```

## Common Tasks

### Parsing command-line flags

Flags are passed as strings in `optsA ...string`:
```go
func parseOpts(opts []string) map[string]string {
    result := make(map[string]string)
    for _, opt := range opts {
        if strings.HasPrefix(opt, "-") {
            parts := strings.SplitN(opt[1:], "=", 2)
            if len(parts) == 2 {
                result[parts[0]] = parts[1]
            }
        }
    }
    return result
}
```
