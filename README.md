# Conedit - Console Text Editor Library

[![Go Reference](https://pkg.go.dev/badge/github.com/topxeq/conedit.svg)](https://pkg.go.dev/github.com/topxeq/conedit)
[![Go Report Card](https://goreportcard.com/badge/github.com/topxeq/conedit)](https://goreportcard.com/report/github.com/topxeq/conedit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**[🌏 中文文档](README_zh.md)**

A lightweight, embeddable terminal text editor library for Go with UTF-8/Chinese support and SSH capability — no CGO required.

## Features

- Full-featured console text editor
- UTF-8 encoding with Chinese character support
- SSH remote file editing
- Auto word wrap (toggleable)
- Keyboard shortcuts for common operations
- Find and replace with regex support
- Undo/redo functionality
- Copy/paste with selection support

## Installation

```bash
go get github.com/topxeq/conedit
```

## Usage as a Library

```go
package main

import (
    "fmt"
    "github.com/topxeq/conedit/editor"
)

func main() {
    // Open editor with default text
    result := editor.ConsoleEditText("Default text content")
    
    // Or open a file
    result = editor.ConsoleEditText("", "-filePath=/path/to/file.txt")
    
    // Or edit remote file via SSH
    result = editor.ConsoleEditText("", 
        "-fromSSH",
        "-sshHost=192.168.1.100",
        "-sshPort=22",
        "-sshUser=root",
        "-sshPass=password",
        "-filePath=/remote/path/file.txt",
    )
    
    // Check result
    if result["status"] == "save" || result["status"] == "saveAs" {
        fmt.Printf("Saved to: %s\n", result["path"])
        fmt.Printf("Content: %s\n", result["text"])
    } else if result["status"] == "cancel" {
        fmt.Println("User cancelled")
    } else if result["status"] == "error" {
        fmt.Printf("Error: %v\n", result["error"])
    }
}
```

## Command Line Usage

Build the command-line editor:

```bash
go build -o conedit ./cmd/editor
```

Usage:

```bash
# Default mode - text input, no file operations
./conedit

# Open file for editing (immediate mode - auto-save on exit)
./conedit file.txt

# Explicit mode selection
./conedit -mode=default        # Text input mode
./conedit -mode=file file.txt  # File mode (returns after save)
./conedit -mode=immediate file.txt  # Immediate mode (auto-save on exit)

# Edit remote file via SSH
./conedit -mode=immediate -fromSSH -sshHost=192.168.1.100 -sshUser=root -filePath=/remote/file.txt

# Show help
./conedit --help
```

### Mode Behavior

| Mode | When | Behavior | Returns |
|------|------|----------|---------|
| `default` | No file args | Text input only | `ok`, `cancel` |
| `file` | With file arg | Edit, return on save | `save`, `saveAs`, `cancel` |
| `immediate` | With file arg | Edit, auto-save on exit | `exit`, `cancel`, `error` |

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Ctrl+S | Save |
| Ctrl+K | Save As |
| Ctrl+X | Exit |
| Ctrl+Q | Force Quit |
| Ctrl+W | Toggle Word Wrap |
| Ctrl+C | Copy |
| Ctrl+V | Paste |
| Ctrl+Z | Undo |
| Ctrl+Y | Redo |
| Ctrl+F | Find (regex supported) |
| Ctrl+H | Replace (regex supported) |
| Ctrl+G | Goto Line |
| Shift+Arrows | Select text |

## Options

| Option | Description |
|--------|-------------|
| `-filePath=PATH` | File path to edit |
| `-mode=MODE` | Editor mode: `default`, `file`, `immediate` (default: `file`) |
| `-fromSSH` | Edit file on SSH server |
| `-sshHost=HOST` | SSH host |
| `-sshPort=PORT` | SSH port (default: 22) |
| `-sshUser=USER` | SSH username |
| `-sshPass=PASS` | SSH password |
| `-sshKeyPath=PATH` | SSH private key path |
| `-mem` | Force in-memory processing (no temp files) |
| `-tmpPath=PATH` | Custom temp directory for large files |

## Return Values

The `ConsoleEditText` function returns a `map[string]interface{}` with the following keys:

| Key | Type | Description |
|-----|------|-------------|
| `text` | string | Current editor content (empty if cancelled or error) |
| `status` | string | One of: `save`, `saveAs`, `cancel`, `error` |
| `path` | string | File path (only when status is `save` or `saveAs`) |
| `error` | string | Error message (only when status is `error`) |

## Requirements

- Go 1.21 or later
- Terminal with UTF-8 support

## Dependencies

- [github.com/gdamore/tcell/v2](https://github.com/gdamore/tcell) - Terminal screen handling
- [golang.org/x/crypto/ssh](https://golang.org/x/crypto) - SSH support

## License

[MIT License](LICENSE)
