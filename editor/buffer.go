package editor

import (
	"regexp"
	"strings"
)

type Buffer struct {
	lines          [][]rune
	undoStack      []Operation
	redoStack      []Operation
	cursorRow      int
	cursorCol      int
	selectionStart *struct {
		row int
		col int
	}
}

type Operation struct {
	kind     string
	row, col int
	oldText  []rune
	newText  []rune
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

	b.undoStack = append(b.undoStack, Operation{
		kind:    "insert",
		row:     row,
		col:     col,
		newText: text,
	})
	b.redoStack = nil
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

	oldText := make([]rune, endCol-col)
	copy(oldText, b.lines[row][col:endCol])

	line := b.lines[row]
	b.lines[row] = append(line[:col], line[endCol:]...)

	b.undoStack = append(b.undoStack, Operation{
		kind:    "delete",
		row:     row,
		col:     col,
		oldText: oldText,
	})
	b.redoStack = nil
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

func (b *Buffer) Undo() {
	if len(b.undoStack) == 0 {
		return
	}
	op := b.undoStack[len(b.undoStack)-1]
	b.undoStack = b.undoStack[:len(b.undoStack)-1]

	switch op.kind {
	case "insert":
		row := op.row
		col := op.col
		count := len(op.newText)
		if col+count <= len(b.lines[row]) {
			b.lines[row] = append(b.lines[row][:col], b.lines[row][col+count:]...)
		}
	case "delete":
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

func (b *Buffer) Find(searchText string, fromRow, fromCol int) (int, int) {
	for row := fromRow; row < len(b.lines); row++ {
		startCol := 0
		if row == fromRow {
			startCol = fromCol
		}
		line := string(b.lines[row])
		idx := strings.Index(line[startCol:], searchText)
		if idx >= 0 {
			return row, startCol + idx
		}
	}
	return -1, -1
}

func (b *Buffer) FindRegex(pattern string, fromRow, fromCol int) (int, int) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return -1, -1
	}
	for row := fromRow; row < len(b.lines); row++ {
		startCol := 0
		if row == fromRow {
			startCol = fromCol
		}
		line := string(b.lines[row])
		loc := re.FindStringIndex(line[startCol:])
		if loc != nil {
			return row, startCol + loc[0]
		}
	}
	return -1, -1
}

func (b *Buffer) ReplaceRegex(pattern, replacement string, fromRow, fromCol int) (int, int, bool) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return -1, -1, false
	}
	for row := fromRow; row < len(b.lines); row++ {
		startCol := 0
		if row == fromRow {
			startCol = fromCol
		}
		line := string(b.lines[row])
		loc := re.FindStringIndex(line[startCol:])
		if loc != nil {
			matchStart := startCol + loc[0]
			matchEnd := startCol + loc[1]
			matchText := line[matchStart:matchEnd]
			newText := re.ReplaceAllString(matchText, replacement)
			b.Delete(row, matchStart, matchEnd-matchStart)
			b.Insert(row, matchStart, []rune(newText))
			return row, matchStart + len([]rune(newText)), true
		}
	}
	return -1, -1, false
}

func (b *Buffer) GetSelectedText() string {
	if b.selectionStart == nil {
		return ""
	}
	startRow, startCol := b.selectionStart.row, b.selectionStart.col
	endRow, endCol := b.cursorRow, b.cursorCol

	if endRow < startRow || (endRow == startRow && endCol < startCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	var result []rune
	for row := startRow; row <= endRow; row++ {
		line := b.lines[row]
		if row == startRow {
			if row == endRow {
				result = append(result, line[startCol:endCol]...)
			} else {
				result = append(result, line[startCol:]...)
				result = append(result, '\n')
			}
		} else if row == endRow {
			result = append(result, line[:endCol]...)
		} else {
			result = append(result, line...)
			result = append(result, '\n')
		}
	}
	return string(result)
}

func (b *Buffer) ClearSelection() {
	b.selectionStart = nil
}

func (b *Buffer) StartSelection() {
	b.selectionStart = &struct {
		row int
		col int
	}{row: b.cursorRow, col: b.cursorCol}
}

func (b *Buffer) DeleteSelection() {
	if b.selectionStart == nil {
		return
	}
	startRow, startCol := b.selectionStart.row, b.selectionStart.col
	endRow, endCol := b.cursorRow, b.cursorCol

	if endRow < startRow || (endRow == startRow && endCol < startCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	if startRow == endRow {
		b.Delete(startRow, startCol, endCol-startCol)
	} else {
		b.lines[startRow] = append(b.lines[startRow][:startCol], b.lines[endRow][endCol:]...)
		b.lines = append(b.lines[:startRow+1], b.lines[endRow+1:]...)
	}
	b.cursorRow = startRow
	b.cursorCol = startCol
	b.selectionStart = nil
}
