package editor

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen         tcell.Screen
	buffer         *Buffer
	opts           map[string]string
	running        bool
	wrap           bool
	clipBoard      string
	status         string
	filePath       string
	inputMode      bool
	inputPrompt    string
	inputBuffer    string
	findRow        int
	findCol        int
	replacePattern string
	mode           string
	unsaved        bool
	// Viewport tracking for performance
	viewportRow    int
	viewportHeight int
}

func ConsoleEditText(defaultTextA string, optsA ...string) map[string]interface{} {
	opts := ParseOpts(optsA)

	text := defaultTextA
	mode := opts["mode"]
	if mode == "" {
		mode = "default"
	}

	filePath := opts["filePath"]
	if filePath == "" {
		for _, opt := range optsA {
			if !strings.HasPrefix(opt, "-") && opt != "" {
				filePath = opt
				break
			}
		}
	}

	if mode == "file" || mode == "immediate" {
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

			if filePath != "" {
				content, err := client.ReadFile(filePath)
				if err != nil {
					return map[string]interface{}{
						"text":   "",
						"status": "error",
						"error":  err.Error(),
					}
				}
				text = content
			}
		} else if filePath != "" {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return map[string]interface{}{
					"text":   "",
					"status": "error",
					"error":  err.Error(),
				}
			}
			text = string(data)
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
		screen:   screen,
		buffer:   NewBuffer(text),
		opts:     opts,
		running:  true,
		wrap:     true,
		status:   "cancel",
		filePath: filePath,
		mode:     mode,
	}

	editor.run()

	switch editor.mode {
	case "default":
		if editor.status == "ok" {
			return map[string]interface{}{
				"text":   editor.buffer.Text(),
				"status": "ok",
			}
		}
		return map[string]interface{}{
			"text":   "",
			"status": "cancel",
		}
	case "immediate":
		if editor.status == "ok" {
			// User exited with Ctrl+X (with or without saving)
			return map[string]interface{}{
				"text":   editor.buffer.Text(),
				"status": "ok",
			}
		} else if editor.status == "error" {
			return map[string]interface{}{
				"text":   "",
				"status": "error",
				"error":  "save failed",
			}
		}
		// Cancel (Ctrl+Q)
		return map[string]interface{}{
			"text":   "",
			"status": "cancel",
		}
	default:
		if editor.status == "error" {
			return map[string]interface{}{
				"text":   "",
				"status": "error",
				"error":  "save failed",
			}
		}
		if editor.status == "cancel" || editor.buffer == nil {
			return map[string]interface{}{
				"text":   "",
				"status": "cancel",
			}
		}
		result := map[string]interface{}{
			"text":   editor.buffer.Text(),
			"status": editor.status,
		}
		if editor.status == "save" || editor.status == "saveAs" {
			result["path"] = editor.filePath
		}
		return result
	}
}

func (e *Editor) run() {
	e.screen.SetStyle(tcell.StyleDefault)
	e.screen.Clear()

	// Initialize viewport
	_, height := e.screen.Size()
	e.viewportHeight = height - 1
	e.viewportRow = 0

	for e.running {
		e.render()
		ev := e.screen.PollEvent()
		if ev == nil {
			break
		}
		e.handleEvent(ev)
	}
}

func (e *Editor) render() {
	if e.buffer == nil {
		return
	}
	e.screen.Clear()

	width, height := e.screen.Size()

	// Ensure viewport is valid
	if e.viewportRow < 0 {
		e.viewportRow = 0
	}
	if e.viewportRow >= e.buffer.LineCount() {
		e.viewportRow = max(0, e.buffer.LineCount()-1)
	}

	// Clamp cursor to valid position
	cursorRow, cursorCol := e.buffer.GetCursor()
	cursorRow, cursorCol = e.buffer.ClampCursor(cursorRow, cursorCol)
	e.buffer.SetCursor(cursorRow, cursorCol)

	// Adjust viewport to keep cursor visible
	e.viewportHeight = height - 1
	if cursorRow < e.viewportRow {
		e.viewportRow = cursorRow
	}
	if cursorRow >= e.viewportRow+e.viewportHeight {
		e.viewportRow = cursorRow - e.viewportHeight + 1
	}

	// Render only visible lines
	screenRow := 0
	for row := e.viewportRow; row < e.buffer.LineCount() && screenRow < e.viewportHeight; row++ {
		line := e.buffer.lines[row]
		lineStr := string(line)

		if e.wrap && width > 0 {
			visualWidth := CalculateVisualWidth(lineStr)
			if visualWidth > width {
				// Handle word wrapping for long lines
				startRune := 0
				visualCol := 0
				runeIdx := 0
				for _, r := range lineStr {
					charWidth := 1
					if isWideRune(r) {
						charWidth = 2
					} else if r == '\t' {
						charWidth = 8
					}
					if visualCol+charWidth > width {
						e.drawLine(screenRow, string(line[startRune:runeIdx]), width)
						screenRow++
						if screenRow >= e.viewportHeight {
							break
						}
						startRune = runeIdx
						visualCol = 0
					}
					visualCol += charWidth
					runeIdx++
				}
				if startRune < len(line) && screenRow < e.viewportHeight {
					e.drawLine(screenRow, string(line[startRune:]), width)
				}
			} else {
				e.drawLine(screenRow, lineStr, width)
			}
		} else {
			e.drawLine(screenRow, lineStr, width)
		}
		screenRow++
	}

	e.drawStatusBar(width, height)
	if e.inputMode {
		e.drawInputMode(width, height)
	}
	e.drawCursor(width, height)
	e.screen.Show()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (e *Editor) drawCursor(width, height int) {
	if e.buffer == nil {
		return
	}
	row, col := e.buffer.GetCursor()

	// Clamp cursor
	row, col = e.buffer.ClampCursor(row, col)

	line := e.buffer.lines[row]
	if col > len(line) {
		col = len(line)
	}

	// Calculate visual column based on character widths
	visualCol := 0
	for i := 0; i < col && i < len(line); i++ {
		r := line[i]
		if r == '\t' {
			visualCol += 8
		} else if isWideRune(r) {
			visualCol += 2
		} else {
			visualCol += 1
		}
	}

	// Calculate screen position relative to viewport
	drawRow := row - e.viewportRow
	drawCol := visualCol

	// Handle word wrap offset
	if e.wrap && width > 0 {
		wrappedLineOffset := 0
		for r := e.viewportRow; r < row; r++ {
			lineStr := string(e.buffer.lines[r])
			visualWidth := CalculateVisualWidth(lineStr)
			if visualWidth > width {
				wrappedLineOffset += (visualWidth - 1) / width
			}
		}
		lineStr := string(line)
		fullVisualWidth := CalculateVisualWidth(lineStr)
		currentLineWrapped := 0
		if fullVisualWidth > width {
			currentLineWrapped = visualCol / width
		}
		drawRow = row - e.viewportRow + wrappedLineOffset + currentLineWrapped
		drawCol = visualCol % width
	}

	// Ensure cursor is within visible area
	if drawRow < 0 || drawRow >= height-1 {
		return // Cursor is outside visible area
	}
	if drawCol >= width {
		drawCol = width - 1
	}

	// Draw cursor
	if drawCol >= 0 {
		if col < len(line) {
			r := line[col]
			e.screen.SetContent(drawCol, drawRow, r, nil, tcell.StyleDefault.Reverse(true))
		} else {
			e.screen.SetContent(drawCol, drawRow, ' ', nil, tcell.StyleDefault.Reverse(true))
		}
	}
}

func (e *Editor) drawLine(row int, text string, width int) {
	col := 0
	for _, r := range text {
		charWidth := 1
		if r == '\t' {
			charWidth = 8
		} else if isWideRune(r) {
			charWidth = 2
		}
		if col+charWidth > width {
			break
		}
		e.screen.SetContent(col, row, r, nil, tcell.StyleDefault)
		col += charWidth
	}
}

func (e *Editor) drawStatusBar(width, height int) {
	row, col := e.buffer.GetCursor()
	cursorInfo := fmt.Sprintf("Line: %d, Col: %d", row+1, col+1)

	// Build status bar content
	var parts []string

	// Show unsaved indicator
	if e.unsaved && e.mode != "default" {
		parts = append(parts, "[MODIFIED]")
	}

	// Show mode-specific help
	var help string
	switch e.mode {
	case "default":
		help = "Ctrl+S:Confirm Ctrl+X:Confirm Ctrl+Q:Cancel"
	case "immediate":
		if e.unsaved {
			help = "Ctrl+S:Save Ctrl+K:SaveAs Ctrl+X:Exit(Confirm)"
		} else {
			help = "Ctrl+S:Save Ctrl+K:SaveAs Ctrl+X:Exit"
		}
	default:
		help = "Ctrl+S:Save Ctrl+K:SaveAs Ctrl+Q:Cancel"
	}
	parts = append(parts, "|")
	parts = append(parts, help)

	// Add cursor info
	parts = append(parts, "|", cursorInfo)

	// Add wrap status
	if e.wrap {
		parts = append(parts, "|", "[Wrap:ON]")
	} else {
		parts = append(parts, "|", "[Wrap:OFF]")
	}

	status := strings.Join(parts, " ")

	// Draw status bar background
	for i := 0; i < width; i++ {
		e.screen.SetContent(i, height-1, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
	}

	// Draw status bar text
	for i, r := range status {
		if i < width {
			e.screen.SetContent(i, height-1, r, nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
		}
	}
}

func (e *Editor) drawInputMode(width, height int) {
	inputLine := height - 2
	prompt := e.inputPrompt + e.inputBuffer

	for i := 0; i < width; i++ {
		e.screen.SetContent(i, inputLine, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack))
	}

	for i, r := range prompt {
		if i < width {
			e.screen.SetContent(i, inputLine, r, nil, tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack))
		}
	}
}

func (e *Editor) handleInputMode(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape {
		if e.inputPrompt == "Save As:" {
			// File mode: cancel and exit
			// Immediate mode: cancel and continue editing
			if e.mode == "file" {
				e.status = "cancel"
				e.running = false
			} else {
				e.inputMode = false
				e.inputBuffer = ""
			}
			return
		}
		if e.inputPrompt == "Save changes? (y/n):" {
			// Cancel exit, continue editing
			e.inputMode = false
			e.inputBuffer = ""
			return
		}
		e.inputMode = false
		e.inputBuffer = ""
		e.inputPrompt = "Find:"
		return
	}
	if ev.Key() == tcell.KeyEnter {
		e.inputMode = false
		switch e.inputPrompt {
		case "Find:":
			pattern := e.inputBuffer
			if pattern != "" {
				row, col := e.buffer.FindRegex(pattern, e.findRow, e.findCol)
				if row >= 0 {
					e.buffer.SetCursor(row, col+len([]rune(pattern)))
					e.findRow = row
					e.findCol = col + len([]rune(pattern))
				}
			}
		case "Replace:":
			e.replacePattern = e.inputBuffer
			e.inputPrompt = "With:"
			e.inputBuffer = ""
			return
		case "With:":
			if e.replacePattern != "" {
				row, col, found := e.buffer.ReplaceRegex(e.replacePattern, e.inputBuffer, e.findRow, e.findCol)
				if found {
					e.buffer.SetCursor(row, col)
					e.findRow = row
					e.findCol = col
					e.unsaved = true
				}
			}
			e.inputPrompt = "Find:"
			e.inputBuffer = ""
		case "Goto:":
			var lineNum int
			fmt.Sscanf(e.inputBuffer, "%d", &lineNum)
			if lineNum > 0 && lineNum <= e.buffer.LineCount() {
				e.buffer.SetCursor(lineNum-1, 0)
			}
		case "Save As:":
			if e.inputBuffer != "" {
				var err error
				if e.opts["fromSSH"] != "" {
					sshConfig := &SSHConfig{
						Host:     e.opts["sshHost"],
						Port:     e.opts["sshPort"],
						User:     e.opts["sshUser"],
						Password: e.opts["sshPass"],
						KeyPath:  e.opts["sshKeyPath"],
					}
					client := NewSSHClient(sshConfig)
					if err = client.Connect(); err != nil {
						e.status = "error"
						if e.mode == "file" {
							e.running = false
						}
						return
					}
					err = client.WriteFile(e.inputBuffer, e.buffer.Text())
					client.Close()
				} else {
					err = os.WriteFile(e.inputBuffer, []byte(e.buffer.Text()), 0644)
				}
				if err != nil {
					e.status = "error"
					if e.mode == "file" {
						e.running = false
					}
				} else {
					e.filePath = e.inputBuffer
					e.unsaved = false
					// File mode: set status and exit after save as
					// Immediate mode: don't change status, continue editing
					if e.mode == "file" {
						e.status = "saveAs"
						e.running = false
					}
				}
			} else {
				// Empty path - cancel
				e.inputMode = false
				e.inputBuffer = ""
				if e.mode == "file" {
					e.status = "cancel"
					e.running = false
				}
				return
			}
		case "Save changes? (y/n):":
			// Immediate mode exit confirmation
			input := strings.ToLower(strings.TrimSpace(e.inputBuffer))
			// Default to "no" (empty input or n) - don't save
			if input == "y" || input == "yes" {
				// Save and exit
				if e.filePath != "" {
					var err error
					if e.opts["fromSSH"] != "" {
						sshConfig := &SSHConfig{
							Host:     e.opts["sshHost"],
							Port:     e.opts["sshPort"],
							User:     e.opts["sshUser"],
							Password: e.opts["sshPass"],
							KeyPath:  e.opts["sshKeyPath"],
						}
						client := NewSSHClient(sshConfig)
						if err = client.Connect(); err != nil {
							e.status = "error"
							e.running = false
							return
						}
						err = client.WriteFile(e.filePath, e.buffer.Text())
						client.Close()
					} else {
						err = os.WriteFile(e.filePath, []byte(e.buffer.Text()), 0644)
					}
					if err != nil {
						e.status = "error"
						e.running = false
						return
					}
				}
			}
			// Either way (y or n/empty), exit with status "ok"
			e.status = "ok"
			e.running = false
		}
		e.inputBuffer = ""
	} else if ev.Key() == tcell.KeyEscape {
		e.inputMode = false
		e.inputBuffer = ""
		e.inputPrompt = "Find:"
	} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if len(e.inputBuffer) > 0 {
			e.inputBuffer = string([]rune(e.inputBuffer)[:len([]rune(e.inputBuffer))-1])
		}
	} else if ev.Key() == tcell.KeyRune {
		e.inputBuffer += string(ev.Rune())
	}
}

func (e *Editor) handleEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if e.inputMode {
			e.handleInputMode(ev)
			return
		}
		if ev.Key() == tcell.KeyCtrlS {
			e.handleCommand(CmdSave)
			return
		}
		if ev.Key() == tcell.KeyCtrlK {
			e.handleCommand(CmdSaveAs)
			return
		}
		if ev.Key() == tcell.KeyCtrlX {
			e.handleCommand(CmdExit)
			return
		}
		if ev.Key() == tcell.KeyCtrlQ {
			e.handleCommand(CmdForceQuit)
			return
		}
		if ev.Key() == tcell.KeyCtrlW {
			e.handleCommand(CmdToggleWrap)
			return
		}
		if ev.Key() == tcell.KeyCtrlC {
			e.handleCommand(CmdCopy)
			return
		}
		if ev.Key() == tcell.KeyCtrlV {
			e.handleCommand(CmdPaste)
			return
		}
		if ev.Key() == tcell.KeyCtrlZ {
			e.handleCommand(CmdUndo)
			return
		}
		if ev.Key() == tcell.KeyCtrlY {
			e.handleCommand(CmdRedo)
			return
		}
		if ev.Key() == tcell.KeyCtrlF {
			e.handleCommand(CmdFind)
			return
		}
		if ev.Key() == tcell.KeyCtrlG {
			e.handleCommand(CmdGotoLine)
			return
		}
		if ev.Key() == tcell.KeyCtrlH {
			e.handleCommand(CmdReplace)
			return
		}

		if ev.Key() == tcell.KeyRune {
			e.buffer.Insert(e.buffer.cursorRow, e.buffer.cursorCol, []rune{ev.Rune()})
			e.buffer.cursorCol++
			e.unsaved = true
		} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
			if e.buffer.cursorCol > 0 {
				e.buffer.cursorCol--
				e.buffer.Delete(e.buffer.cursorRow, e.buffer.cursorCol, 1)
				e.unsaved = true
			}
		} else if ev.Key() == tcell.KeyEnter {
			newLine := make([]rune, 0)
			if e.buffer.cursorCol < len(e.buffer.lines[e.buffer.cursorRow]) {
				newLine = append(newLine, e.buffer.lines[e.buffer.cursorRow][e.buffer.cursorCol:]...)
				e.buffer.lines[e.buffer.cursorRow] = e.buffer.lines[e.buffer.cursorRow][:e.buffer.cursorCol]
			}
			e.buffer.lines = append(e.buffer.lines[:e.buffer.cursorRow+1], append([][]rune{newLine}, e.buffer.lines[e.buffer.cursorRow+1:]...)...)
			e.buffer.cursorRow++
			e.buffer.cursorCol = 0
			e.unsaved = true
		} else if ev.Key() == tcell.KeyUp {
			e.buffer.ClearSelection()
			if e.buffer.cursorRow > 0 {
				e.buffer.cursorRow--
				if e.buffer.cursorCol > e.buffer.LineLen(e.buffer.cursorRow) {
					e.buffer.cursorCol = e.buffer.LineLen(e.buffer.cursorRow)
				}
			}
		} else if ev.Key() == tcell.KeyDown {
			e.buffer.ClearSelection()
			if e.buffer.cursorRow < e.buffer.LineCount()-1 {
				e.buffer.cursorRow++
				if e.buffer.cursorCol > e.buffer.LineLen(e.buffer.cursorRow) {
					e.buffer.cursorCol = e.buffer.LineLen(e.buffer.cursorRow)
				}
			}
		} else if ev.Key() == tcell.KeyLeft {
			e.buffer.ClearSelection()
			if e.buffer.cursorCol > 0 {
				e.buffer.cursorCol--
			}
		} else if ev.Key() == tcell.KeyRight {
			e.buffer.ClearSelection()
			if e.buffer.cursorCol < e.buffer.LineLen(e.buffer.cursorRow) {
				e.buffer.cursorCol++
			}
		} else if ev.Modifiers()&tcell.ModShift != 0 {
			if e.buffer.selectionStart == nil {
				e.buffer.StartSelection()
			}
			if ev.Key() == tcell.KeyUp && e.buffer.cursorRow > 0 {
				e.buffer.cursorRow--
				if e.buffer.cursorCol > e.buffer.LineLen(e.buffer.cursorRow) {
					e.buffer.cursorCol = e.buffer.LineLen(e.buffer.cursorRow)
				}
			} else if ev.Key() == tcell.KeyDown && e.buffer.cursorRow < e.buffer.LineCount()-1 {
				e.buffer.cursorRow++
				if e.buffer.cursorCol > e.buffer.LineLen(e.buffer.cursorRow) {
					e.buffer.cursorCol = e.buffer.LineLen(e.buffer.cursorRow)
				}
			} else if ev.Key() == tcell.KeyLeft && e.buffer.cursorCol > 0 {
				e.buffer.cursorCol--
			} else if ev.Key() == tcell.KeyRight && e.buffer.cursorCol < e.buffer.LineLen(e.buffer.cursorRow) {
				e.buffer.cursorCol++
			}
		} else if ev.Key() == tcell.KeyHome {
			e.buffer.cursorCol = 0
		} else if ev.Key() == tcell.KeyEnd {
			e.buffer.cursorCol = e.buffer.LineLen(e.buffer.cursorRow)
		} else if ev.Key() == tcell.KeyPgUp {
			if e.buffer.cursorRow > 0 {
				e.buffer.cursorRow--
			}
		} else if ev.Key() == tcell.KeyPgDn {
			if e.buffer.cursorRow < e.buffer.LineCount()-1 {
				e.buffer.cursorRow++
			}
		}
	}
}

func (e *Editor) handleCommand(cmd Command) {
	switch cmd {
	case CmdSave:
		if e.mode == "default" {
			e.status = "ok"
			e.running = false
			return
		}

		if e.filePath != "" {
			var err error
			if e.opts["fromSSH"] != "" {
				sshConfig := &SSHConfig{
					Host:     e.opts["sshHost"],
					Port:     e.opts["sshPort"],
					User:     e.opts["sshUser"],
					Password: e.opts["sshPass"],
					KeyPath:  e.opts["sshKeyPath"],
				}
				client := NewSSHClient(sshConfig)
				if err := client.Connect(); err != nil {
					e.status = "error"
					if e.mode == "file" {
						e.running = false
					}
					return
				}
				err = client.WriteFile(e.filePath, e.buffer.Text())
				client.Close()
			} else {
				err = os.WriteFile(e.filePath, []byte(e.buffer.Text()), 0644)
			}
			if err != nil {
				e.status = "error"
				if e.mode == "file" {
					e.running = false
				}
				return
			}
			e.unsaved = false

			// File mode: set status and exit after save
			// Immediate mode: don't change status, continue editing
			if e.mode == "file" {
				e.status = "save"
				e.running = false
			}
		} else {
			e.inputMode = true
			e.inputPrompt = "Save As:"
			e.inputBuffer = ""
		}
	case CmdSaveAs:
		if e.mode == "default" {
			e.status = "ok"
			e.running = false
			return
		}
		// Prompt for file path
		// File mode: exit after save as is completed
		// Immediate mode: continue editing after save as
		e.inputMode = true
		e.inputPrompt = "Save As:"
		e.inputBuffer = ""
	case CmdCopy:
		selected := e.buffer.GetSelectedText()
		if selected != "" {
			e.clipBoard = selected
			e.buffer.DeleteSelection()
			e.unsaved = true
		}
	case CmdPaste:
		if e.clipBoard != "" {
			if e.buffer.selectionStart != nil {
				e.buffer.DeleteSelection()
			}
			e.buffer.Insert(e.buffer.cursorRow, e.buffer.cursorCol, []rune(e.clipBoard))
			e.buffer.cursorCol += len([]rune(e.clipBoard))
			e.unsaved = true
		}
	case CmdUndo:
		e.buffer.Undo()
		e.unsaved = true
	case CmdRedo:
		e.buffer.Redo()
		e.unsaved = true
	case CmdToggleWrap:
		e.wrap = !e.wrap
	case CmdExit:
		if e.mode == "default" {
			e.status = "ok"
			e.running = false
			return
		}

		// Immediate mode: check if modified, prompt to save if needed
		if e.mode == "immediate" {
			if e.unsaved && e.filePath != "" {
				// Prompt user to confirm save
				e.inputMode = true
				e.inputPrompt = "Save changes? (y/n):"
				e.inputBuffer = ""
				return
			}
			// No changes, exit directly
			e.status = "ok"
			e.running = false
			return
		}

		// File mode: no Ctrl+X handling (only Ctrl+Q for cancel)
		e.buffer = nil
		e.status = "cancel"
		e.running = false
	case CmdForceQuit:
		// Force quit without saving in any mode
		e.buffer = nil
		e.status = "cancel"
		e.running = false
	case CmdFind:
		e.inputMode = true
		e.inputPrompt = "Find:"
		e.inputBuffer = ""
		e.findRow = 0
		e.findCol = 0
	case CmdReplace:
		e.inputMode = true
		e.inputPrompt = "Replace:"
		e.inputBuffer = ""
		e.replacePattern = ""
		e.findRow = 0
		e.findCol = 0
	case CmdGotoLine:
		e.inputMode = true
		e.inputPrompt = "Goto:"
		e.inputBuffer = ""
	}
}
