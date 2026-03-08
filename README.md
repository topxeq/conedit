# Conedit - Console Text Editor Library

A lightweight command-line text editor library for Go with SSH support.

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
# Edit a file
./conedit file.txt

# Edit with options
./conedit -filePath=file.txt

# Edit remote file via SSH
./conedit -fromSSH -sshHost=192.168.1.100 -sshUser=root -filePath=/remote/file.txt

# Show help
./conedit --help
```

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

## Modes

| Mode | Description | Return Status |
|------|-------------|---------------|
| `default` | Simple text input, no file operations | `ok`, `cancel` |
| `file` | Open/edit file, return after save action | `save`, `saveAs`, `cancel`, `error` |
| `immediate` | Auto-save on exit, direct file operations | `exit`, `cancel`, `error` |

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

MIT License
