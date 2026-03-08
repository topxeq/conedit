# Console Text Editor Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a console-based text editor with SSH support in Go

**Architecture:** Single package with modular components (buffer, screen, input, command, sshclient)

**Tech Stack:** Go, tcell, golang.org/x/crypto/ssh

---

### Task 1: Initialize project and create go.mod

**Files:**
- Create: `/mnt1/aiprjs/conedit/go.mod`
- Create: `/mnt1/aiprjs/conedit/main.go`

**Step 1: Create go.mod**

```go
module console-editor

go 1.21

require (
    github.com/gdamore/tcell/v2 v2.6.0
    golang.org/x/crypto v0.14.0
)
```

**Step 2: Create main.go skeleton**

```go
package main

import "fmt"

func consoleEditText(defaultTextA string, optsA ...string) map[string]interface{} {
    return map[string]interface{}{
        "text":   "",
        "status": "cancel",
    }
}

func main() {
    result := consoleEditText("")
    fmt.Println(result)
}
```

**Step 3: Run go mod tidy**

Run: `cd /mnt1/aiprjs/conedit && go mod tidy`

---

### Task 2: Create text buffer with line management

**Files:**
- Create: `/mnt1/aiprjs/conedit/buffer.go`

**Step 1: Write test**

```go
package main

import "testing"

func TestBuffer_NewBuffer(t *testing.T) {
    buf := NewBuffer("hello\nworld")
    if len(buf.lines) != 2 {
        t.Errorf("expected 2 lines, got %d", len(buf.lines))
    }
}

func TestBuffer_Insert(t *testing.T) {
    buf := NewBuffer("")
    buf.Insert(0, 0, []rune("hello"))
    if string(buf.lines[0]) != "hello" {
        t.Errorf("expected 'hello', got '%s'", string(buf.lines[0]))
    }
}

func TestBuffer_Delete(t *testing.T) {
    buf := NewBuffer("hello")
    buf.Delete(0, 0, 5)
    if len(buf.lines[0]) != 0 {
        t.Errorf("expected empty line, got '%s'", string(buf.lines[0]))
    }
}
```

**Step 2: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestBuffer`
Expected: FAIL - undefined

**Step 3: Implement Buffer**

```go
package main

import "unicode/utf8"

type Buffer struct {
    lines     [][]rune
    undoStack []Operation
    redoStack []Operation
    cursorRow int
    cursorCol int
}

type Operation struct {
    kind      string
    row, col  int
    oldText   []rune
    newText   []rune
}

func NewBuffer(text string) *Buffer {
    lines := [][]rune{{}}
    for _, r := range text {
        if r == '\n' {
            lines = append(lines, []rune{})
        } else {
            lines[len(lines)-1] = append(lines[len(lines)-1], r)
        }
    }
    return &Buffer{lines: lines}
}

func (b *Buffer) Text() string {
    var result []rune
    for i, line := range b.lines {
        result = append(result, line...)
        if i < len(b.lines)-1 {
            result = append(result, '\n')
        }
    }
    return string(result)
}

func (b *Buffer) Insert(row, col int, text []rune) {
    if row < 0 || row >= len(b.lines) {
        return
    }
    if col < 0 {
        col = 0
    }
    if col > len(b.lines[row]) {
        col = len(b.lines[row])
    }
    
    line := b.lines[row]
    newLine := make([]rune, 0, len(line)+len(text))
    newLine = append(newLine, line[:col]...)
    newLine = append(newLine, text...)
    newLine = append(newLine, line[col:]...)
    b.lines[row] = newLine
}

func (b *Buffer) Delete(row, col, count int) {
    if row < 0 || row >= len(b.lines) {
        return
    }
    if col < 0 || col >= len(b.lines[row]) {
        return
    }
    endCol := col + count
    if endCol > len(b.lines[row]) {
        endCol = len(b.lines[row])
    }
    
    line := b.lines[row]
    b.lines[row] = append(line[:col], line[endCol:]...)
}

func (b *Buffer) LineCount() int {
    return len(b.lines)
}

func (b *Buffer) LineLen(row int) int {
    if row < 0 || row >= len(b.lines) {
        return 0
    }
    return len(b.lines[row])
}

func (b *Buffer) GetCursor() (int, int) {
    return b.cursorRow, b.cursorCol
}

func (b *Buffer) SetCursor(row, col int) {
    if row < 0 {
        row = 0
    }
    if row >= len(b.lines) {
        row = len(b.lines) - 1
    }
    if col < 0 {
        col = 0
    }
    if col > len(b.lines[row]) {
        col = len(b.lines[row])
    }
    b.cursorRow = row
    b.cursorCol = col
}
```

**Step 4: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestBuffer`
Expected: PASS

---

### Task 3: Add undo/redo support

**Files:**
- Modify: `/mnt1/aiprjs/conedit/buffer.go`

**Step 1: Write test**

```go
func TestBuffer_Undo(t *testing.T) {
    buf := NewBuffer("")
    buf.Insert(0, 0, []rune("hello"))
    buf.Undo()
    if len(buf.lines[0]) != 0 {
        t.Errorf("expected empty after undo, got '%s'", string(buf.lines[0]))
    }
}

func TestBuffer_Redo(t *testing.T) {
    buf := NewBuffer("")
    buf.Insert(0, 0, []rune("hello"))
    buf.Undo()
    buf.Redo()
    if string(buf.lines[0]) != "hello" {
        t.Errorf("expected 'hello' after redo, got '%s'", string(buf.lines[0]))
    }
}
```

**Step 2: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestBuffer_Undo`
Expected: FAIL - undefined Undo/Redo

**Step 3: Implement Undo/Redo**

Add to buffer.go:

```go
func (b *Buffer) Undo() {
    if len(b.undoStack) == 0 {
        return
    }
    op := b.undoStack[len(b.undoStack)-1]
    b.undoStack = b.undoStack[:len(b.undoStack)-1]
    
    switch op.kind {
    case "insert":
        // Revert insert by deleting
        row := op.row
        col := op.col
        count := len(op.newText)
        if col+count <= len(b.lines[row]) {
            b.lines[row] = append(b.lines[row][:col], b.lines[row][col+count:]...)
        }
    case "delete":
        // Revert delete by inserting
        b.Insert(op.row, op.col, op.oldText)
    }
    
    b.redoStack = append(b.redoStack, op)
}

func (b *Buffer) Redo() {
    if len(b.redoStack) == 0 {
        return
    }
    op := b.redoStack[len(b.redoStack)-1]
    b.redoStack = b.redoStack[:len(b.redoStack)-1]
    
    switch op.kind {
    case "insert":
        b.Insert(op.row, op.col, op.newText)
    case "delete":
        b.Delete(op.row, op.col, len(op.oldText))
    }
    
    b.undoStack = append(b.undoStack, op)
}
```

**Step 4: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestBuffer_Undo`
Expected: PASS

---

### Task 4: Create screen rendering

**Files:**
- Create: `/mnt1/aiprjs/conedit/screen.go`

**Step 1: Write test**

```go
func TestScreen_CalculateVisualWidth(t *testing.T) {
    tests := []struct {
        input    string
        expected int
    }{
        {"hello", 5},
        {"你好", 2},
        {"hi你好", 4},
    }
    for _, tt := range tests {
        result := CalculateVisualWidth(tt.input)
        if result != tt.expected {
            t.Errorf("CalculateVisualWidth(%q) = %d, want %d", tt.input, result, tt.expected)
        }
    }
}
```

**Step 2: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestScreen`
Expected: FAIL

**Step 3: Implement screen utilities**

```go
package main

import (
    "unicode/utf8"
)

func CalculateVisualWidth(s string) int {
    width := 0
    for _, r := range s {
        if r == '\t' {
            width += 8
        } else if r >= 0x4E00 && r <= 0x9FFF { // CJK
            width += 2
        } else if r >= 0xAC00 && r <= 0xD7AF { // Hangul
            width += 2
        } else if r >= 0x3000 && r <= 0x303F { // CJK Symbols
            width += 2
        } else if r >= 0xFF00 && r <= 0xFFEF { // Fullwidth
            width += 2
        } else {
            width += 1
        }
    }
    return width
}

func RuneToVisualIndex(s string, visualPos int) int {
    width := 0
    for i := range s {
        r, _ := utf8.DecodeRuneInString(s[i:])
        charWidth := 1
        if r == '\t' {
            charWidth = 8
        } else if r >= 0x4E00 && r <= 0x9FFF {
            charWidth = 2
        } else if r >= 0xAC00 && r <= 0xD7AF {
            charWidth = 2
        } else if r >= 0x3000 && r <= 0x303F {
            charWidth = 2
        } else if r >= 0xFF00 && r <= 0xFFEF {
            charWidth = 2
        }
        if width + charWidth > visualPos {
            return i
        }
        width += charWidth
    }
    return len(s)
}
```

**Step 4: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestScreen`
Expected: PASS

---

### Task 5: Create input handling

**Files:**
- Create: `/mnt1/aiprjs/conedit/input.go`

**Step 1: Write test**

```go
func TestParseOpts(t *testing.T) {
    opts := []string{
        "-filePath=/path/to/file.txt",
        "-sshHost=192.168.1.100",
        "-sshPort=22",
        "-sshUser=root",
        "-sshPass=abc123",
        "-tmpPath=/tmp",
        "-mem",
        "-fromSSH",
    }
    result := ParseOpts(opts)
    
    if result["filePath"] != "/path/to/file.txt" {
        t.Errorf("filePath = %s, want /path/to/file.txt", result["filePath"])
    }
    if result["sshHost"] != "192.168.1.100" {
        t.Errorf("sshHost = %s, want 192.168.1.100", result["sshHost"])
    }
    if result["mem"] != "" {
        t.Errorf("mem = %s, want empty string", result["mem"])
    }
    if result["fromSSH"] != "" {
        t.Errorf("fromSSH = %s, want empty string", result["fromSSH"])
    }
}
```

**Step 2: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestParseOpts`
Expected: FAIL

**Step 3: Implement ParseOpts**

```go
package main

import "strings"

func ParseOpts(opts []string) map[string]string {
    result := make(map[string]string)
    for _, opt := range opts {
        if strings.HasPrefix(opt, "-") {
            key := strings.TrimPrefix(opt, "-")
            if idx := strings.Index(key, "="); idx > 0 {
                result[key[:idx]] = key[idx+1:]
            } else {
                result[key] = ""
            }
        }
    }
    return result
}
```

**Step 4: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestParseOpts`
Expected: PASS

---

### Task 6: Create SSH client

**Files:**
- Create: `/mnt1/aiprjs/conedit/sshclient.go`

**Step 1: Write test**

```go
func TestSSHClient_Connect(t *testing.T) {
    // Skip if no SSH server available
    t.Skip("Requires SSH server")
}
```

**Step 2: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestSSHClient`
Expected: SKIP

**Step 3: Implement SSH client**

```go
package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    "golang.org/x/crypto/ssh"
)

type SSHConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    KeyPath  string
}

type SSHClient struct {
    config *SSHConfig
    client *ssh.Client
}

func NewSSHClient(config *SSHConfig) *SSHClient {
    return &SSHClient{config: config}
}

func (s *SSHClient) Connect() error {
    authMethods := []ssh.AuthMethod{}
    
    if s.config.Password != "" {
        authMethods = append(authMethods, ssh.Password(s.config.Password))
    }
    
    if s.config.KeyPath != "" {
        key, err := ioutil.ReadFile(s.config.KeyPath)
        if err == nil {
            signer, err := ssh.ParsePrivateKey(key)
            if err == nil {
                authMethods = append(authMethods, ssh.PublicKeys(signer))
            }
        }
    }
    
    if len(authMethods) == 0 {
        return fmt.Errorf("no authentication method available")
    }
    
    config := &ssh.ClientConfig{
        User: s.config.User,
        Auth: authMethods,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    
    addr := s.config.Host
    if s.config.Port != "" {
        addr = addr + ":" + s.config.Port
    } else {
        addr = addr + ":22"
    }
    
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }
    
    s.client = client
    return nil
}

func (s *SSHClient) ReadFile(path string) (string, error) {
    if s.client == nil {
        return "", fmt.Errorf("not connected")
    }
    
    session, err := s.client.NewSession()
    if err != nil {
        return "", err
    }
    defer session.Close()
    
    output, err := session.CombinedOutput("cat " + path)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %w", err)
    }
    
    return string(output), nil
}

func (s *SSHClient) WriteFile(path, content string) error {
    if s.client == nil {
        return fmt.Errorf("not connected")
    }
    
    session, err := s.client.NewSession()
    if err != nil {
        return err
    }
    defer session.Close()
    
    session.Stdin = strings.NewReader(content)
    err = session.Run("cat > " + path)
    if err != nil {
        return fmt.Errorf("failed to write file: %w", err)
    }
    
    return nil
}

func (s *SSHClient) Close() {
    if s.client != nil {
        s.client.Close()
    }
}
```

**Step 4: Add import**

Add `"strings"` import to sshclient.go

**Step 5: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestSSHClient`
Expected: SKIP

---

### Task 7: Create command handling

**Files:**
- Create: `/mnt1/aiprjs/conedit/command.go`

**Step 1: Write test**

```go
func TestCommand_Parse(t *testing.T) {
    tests := []struct {
        input    string
        expected Command
    }{
        {"\x13", CmdSave},        // Ctrl+S
        {"\x0b", CmdSaveAs},      // Ctrl+K
        {"\x03", CmdCopy},        // Ctrl+C
        {"\x16", CmdPaste},       // Ctrl+V
        {"\x1a", CmdUndo},        // Ctrl+Z
        {"\x19", CmdRedo},        // Ctrl+Y
        {"\x06", CmdFind},        // Ctrl+F
        {"\x08", CmdReplace},     // Ctrl+H
        {"\x07", CmdGotoLine},    // Ctrl+G
        {"\x18", CmdExit},        // Ctrl+X
        {"\x11", CmdForceQuit},   // Ctrl+Q
        {"\x17", CmdToggleWrap},  // Ctrl+W
    }
    for _, tt := range tests {
        result := ParseCommand(tt.input)
        if result != tt.expected {
            t.Errorf("ParseCommand(%q) = %v, want %v", tt.input, result, tt.expected)
        }
    }
}
```

**Step 2: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestCommand`
Expected: FAIL

**Step 3: Implement command handling**

```go
package main

type Command int

const (
    CmdNone Command = iota
    CmdSave
    CmdSaveAs
    CmdCopy
    CmdPaste
    CmdUndo
    CmdRedo
    CmdFind
    CmdReplace
    CmdGotoLine
    CmdExit
    CmdForceQuit
    CmdToggleWrap
)

func ParseCommand(input string) Command {
    if len(input) == 0 {
        return CmdNone
    }
    r := rune(input[0])
    switch r {
    case 0x13: // Ctrl+S
        return CmdSave
    case 0x0B: // Ctrl+K
        return CmdSaveAs
    case 0x03: // Ctrl+C
        return CmdCopy
    case 0x16: // Ctrl+V
        return CmdPaste
    case 0x1A: // Ctrl+Z
        return CmdUndo
    case 0x19: // Ctrl+Y
        return CmdRedo
    case 0x06: // Ctrl+F
        return CmdFind
    case 0x08: // Ctrl+H
        return CmdReplace
    case 0x07: // Ctrl+G
        return CmdGotoLine
    case 0x18: // Ctrl+X
        return CmdExit
    case 0x11: // Ctrl+Q
        return CmdForceQuit
    case 0x17: // Ctrl+W
        return CmdToggleWrap
    }
    return CmdNone
}
```

**Step 4: Run test**

Run: `cd /mnt1/aiprjs/conedit && go test -v -run TestCommand`
Expected: PASS

---

### Task 8: Create main editor logic

**Files:**
- Modify: `/mnt1/aiprjs/conedit/main.go`

**Step 1: Implement consoleEditText function**

```go
package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/gdamore/tcell/v2"
)

type Editor struct {
    screen    tcell.Screen
    buffer    *Buffer
    opts      map[string]string
    running   bool
    wrap      bool
    clipBoard string
}

func consoleEditText(defaultTextA string, optsA ...string) map[string]interface{} {
    opts := ParseOpts(optsA)
    
    text := defaultTextA
    
    if opts["fromSSH"] != "" {
        sshConfig := &SSHConfig{
            Host:     opts["sshHost"],
            Port:     opts["sshPort"],
            User:     opts["sshUser"],
            Password: opts["sshPass"],
            KeyPath:  opts["sshKeyPath"],
        }
        client := NewSSHClient(sshConfig)
        if err := client.Connect(); err != nil {
            return map[string]interface{}{
                "text":   "",
                "status": "error",
                "error":  err.Error(),
            }
        }
        defer client.Close()
        
        if opts["filePath"] != "" {
            content, err := client.ReadFile(opts["filePath"])
            if err != nil {
                return map[string]interface{}{
                    "text":   "",
                    "status": "error",
                    "error":  err.Error(),
                }
            }
            text = content
        }
    } else if opts["filePath"] != "" {
        data, err := os.ReadFile(opts["filePath"])
        if err == nil {
            text = string(data)
        }
    }
    
    if opts["mem"] == "" && len(text) > 10*1024*1024 {
        tmpPath := opts["tmpPath"]
        if tmpPath == "" {
            tmpPath = os.TempDir()
        }
        tmpFile, err := os.CreateTemp(tmpPath, "editor_*.tmp")
        if err != nil {
            return map[string]interface{}{
                "text":   "",
                "status": "error",
                "error":  err.Error(),
            }
        }
        tmpFile.WriteString(text)
        tmpFile.Close()
        defer os.Remove(tmpFile.Name())
    }
    
    screen, err := tcell.NewScreen()
    if err != nil {
        return map[string]interface{}{
            "text":   "",
            "status": "error",
            "error":  err.Error(),
        }
    }
    defer screen.Fini()
    
    if err := screen.Init(); err != nil {
        return map[string]interface{}{
            "text":   "",
            "status": "error",
            "error":  err.Error(),
        }
    }
    
    editor := &Editor{
        screen: screen,
        buffer: NewBuffer(text),
        opts:   opts,
        running: true,
        wrap:   true,
    }
    
    editor.run()
    
    if editor.buffer == nil {
        return map[string]interface{}{
            "text":   "",
            "status": "cancel",
        }
    }
    
    return map[string]interface{}{
        "text":   editor.buffer.Text(),
        "status": "save",
    }
}

func (e *Editor) run() {
    e.screen.SetStyle(tcell.StyleDefault)
    e.screen.Clear()
    
    for e.running {
        e.render()
        ev := e.screen.PollEvent()
        e.handleEvent(ev)
    }
}

func (e *Editor) render() {
    e.screen.Clear()
    
    width, height := e.screen.Size()
    
    for row := 0; row < height-1 && row < e.buffer.LineCount(); row++ {
        line := e.buffer.lines[row]
        lineStr := string(line)
        
        if e.wrap && width > 0 {
            visualWidth := CalculateVisualWidth(lineStr)
            if visualWidth > width {
                start := 0
                col := 0
                for i, r := range lineStr {
                    charWidth := 1
                    if r >= 0x4E00 && r <= 0x9FFF || r >= 0xAC00 && r <= 0xD7AF {
                        charWidth = 2
                    }
                    if col+charWidth > width {
                        e.drawLine(row-start, string(line[start:i]), width)
                        row++
                        if row >= height-1 {
                            break
                        }
                        start = i
                        col = 0
                    }
                    col += charWidth
                }
                if start < len(lineStr) {
                    e.drawLine(row-start, string(line[start:]), width)
                }
            } else {
                e.drawLine(row, lineStr, width)
            }
        } else {
            e.drawLine(row, lineStr, width)
        }
    }
    
    e.drawStatusBar(width, height)
    e.screen.Show()
}

func (e *Editor) drawLine(row int, text string, width int) {
    for i, r := range text {
        e.screen.SetContent(i, row, r, nil, tcell.StyleDefault)
    }
}

func (e *Editor) drawStatusBar(width, height int) {
    status := "Ctrl+S:Save Ctrl+K:SaveAs Ctrl+X:Exit"
    if e.wrap {
        status += " [Wrap:ON]"
    } else {
        status += " [Wrap:OFF]"
    }
    
    for i := 0; i < width; i++ {
        e.screen.SetContent(i, height-1, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
    }
    
    for i, r := range status {
        if i < width {
            e.screen.SetContent(i, height-1, r, nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
        }
    }
}

func (e *Editor) handleEvent(ev tcell.Event) {
    switch ev := ev.(type) {
    case *tcell.EventKey:
        cmd := ParseCommand(string(ev.Key()))
        if cmd != CmdNone {
            e.handleCommand(cmd)
            return
        }
        
        if ev.Key() == tcell.KeyRune {
            e.buffer.Insert(e.buffer.cursorRow, e.buffer.cursorCol, []rune{ev.Rune()})
            e.buffer.cursorCol++
        } else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
            if e.buffer.cursorCol > 0 {
                e.buffer.cursorCol--
                e.buffer.Delete(e.buffer.cursorRow, e.buffer.cursorCol, 1)
            }
        } else if ev.Key() == tcell.KeyEnter {
            e.buffer.lines = append(e.buffer.lines[:e.buffer.cursorRow+1], append([][]rune{{}}, e.buffer.lines[e.buffer.cursorRow+1:]...)...)
            e.buffer.cursorRow++
            e.buffer.cursorCol = 0
        } else if ev.Key() == tcell.KeyUp {
            if e.buffer.cursorRow > 0 {
                e.buffer.cursorRow--
            }
        } else if ev.Key() == tcell.KeyDown {
            if e.buffer.cursorRow < e.buffer.LineCount()-1 {
                e.buffer.cursorRow++
            }
        } else if ev.Key() == tcell.KeyLeft {
            if e.buffer.cursorCol > 0 {
                e.buffer.cursorCol--
            }
        } else if ev.Key() == tcell.KeyRight {
            if e.buffer.cursorCol < e.buffer.LineLen(e.buffer.cursorRow) {
                e.buffer.cursorCol++
            }
        } else if ev.Key() == tcell.KeyCtrlC {
            e.handleCommand(CmdCopy)
        } else if ev.Key() == tcell.KeyCtrlV {
            e.handleCommand(CmdPaste)
        } else if ev.Key() == tcell.KeyCtrlZ {
            e.handleCommand(CmdUndo)
        } else if ev.Key() == tcell.KeyCtrlY {
            e.handleCommand(CmdRedo)
        } else if ev.Key() == tcell.KeyCtrlW {
            e.handleCommand(CmdToggleWrap)
        }
    }
}

func (e *Editor) handleCommand(cmd Command) {
    switch cmd {
    case CmdSave:
        if e.opts["filePath"] != "" {
            err := os.WriteFile(e.opts["filePath"], []byte(e.buffer.Text()), 0644)
            if err != nil {
                return
            }
        }
        e.running = false
    case CmdSaveAs:
        // Simplified: just close as save
        e.running = false
    case CmdCopy:
        row, col := e.buffer.GetCursor()
        if col < len(e.buffer.lines[row]) {
            e.clipBoard = string(e.buffer.lines[row][col])
        }
    case CmdPaste:
        if e.clipBoard != "" {
            e.buffer.Insert(e.buffer.cursorRow, e.buffer.cursorCol, []rune(e.clipBoard))
            e.buffer.cursorCol += len(e.clipBoard)
        }
    case CmdUndo:
        e.buffer.Undo()
    case CmdRedo:
        e.buffer.Redo()
    case CmdToggleWrap:
        e.wrap = !e.wrap
    case CmdExit:
        e.buffer = nil
        e.running = false
    case CmdForceQuit:
        e.buffer = nil
        e.running = false
    }
}

func main() {
    result := consoleEditText("Hello 世界\n第二行")
    fmt.Println(result)
}
```

**Step 2: Run go build**

Run: `cd /mnt1/aiprjs/conedit && go build -o console_editor .`
Expected: SUCCESS

**Step 3: Run the application**

Run: `cd /mnt1/aiprjs/conedit && ./console_editor`

---

### Task 9: Add find/replace functionality

**Files:**
- Modify: `/mnt1/aiprjs/conedit/command.go`

**Step 1: Add find command**

Add to command.go:

```go
const (
    // ... existing commands
    CmdFindNext
    CmdFindPrev
)
```

**Step 2: Modify editor to support find**

Update editor.go to add find state and dialog

---

### Task 10: Final testing and validation

**Step 1: Run all tests**

Run: `cd /mnt1/aiprjs/conedit && go test -v ./...`

**Step 2: Run linting**

Run: `cd /mnt1/aiprjs/conedit && go vet && go fmt ./...`

**Step 3: Build final binary**

Run: `cd /mnt1/aiprjs/conedit && go build -o console_editor .`
