# Console Text Editor Design

## Project Overview

Build a console-based text editor in Go with SSH support.

**Main Function:**
```go
func consoleEditText(defaultTextA string, optsA ...string) map[string]interface{}
```

## Architecture

### Single Package Structure
```
.
├── main.go         # Entry point
├── editor.go       # Core editor logic
├── buffer.go       # Text buffer management
├── screen.go       # Screen rendering
├── input.go        # Input handling
├── command.go      # Command handling
├── sshclient.go    # SSH client
└── go.mod          # Go module
```

### Dependencies
- `github.com/gdamore/tcell` - Terminal UI
- `golang.org/x/crypto/ssh` - SSH support

## Components

### Editor
- `Buffer` - Text buffer with undo/redo support
- `Cursor` - Cursor position (rune-based index)
- `View` - View area (line numbers, edit area)
- `Screen` - tcell screen wrapper

### TextBuffer
- Store as `[]rune` (UTF-8 support)
- `lines [][]rune` - Line-by-line storage
- `undoStack`, `redoStack` - Undo/redo stacks

### Return Format
```go
map[string]interface{}{
    "text":   "",        // Empty on save
    "status": "save|cancel|error",
    "error":  "",        // Only on error
}
```

## Key Features

- Chinese support: All operations use rune, not byte
- Cursor movement: Use `utf8.RuneCountInString` for offset
- Auto-wrap: Calculate screen width and character display width
- Temp files: >10MB use temp file, <10MB in-memory

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Ctrl+S   | Save |
| Ctrl+K   | Save As |
| Ctrl+C   | Copy |
| Ctrl+V   | Paste |
| Ctrl+Z   | Undo |
| Ctrl+Y   | Redo |
| Ctrl+F   | Find |
| Ctrl+H   | Replace |
| Ctrl+G   | Goto Line |
| Ctrl+X   | Exit |
| Ctrl+Q   | Force Quit |
| Ctrl+W   | Toggle Wrap |
