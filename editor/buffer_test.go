package editor

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

func TestBuffer_Find(t *testing.T) {
	buf := NewBuffer("hello world\ntest line")
	row, col := buf.Find("world", 0, 0)
	if row != 0 || col != 6 {
		t.Errorf("Find returned (%d, %d), want (0, 6)", row, col)
	}
}

func TestBuffer_FindNotFound(t *testing.T) {
	buf := NewBuffer("hello world\ntest line")
	row, col := buf.Find("xyz", 0, 0)
	if row != -1 || col != -1 {
		t.Errorf("Find returned (%d, %d), want (-1, -1)", row, col)
	}
}

func TestBuffer_FindFromMiddle(t *testing.T) {
	buf := NewBuffer("hello world\nworld again")
	row, col := buf.Find("world", 0, 8)
	if row != 1 || col != 0 {
		t.Errorf("Find returned (%d, %d), want (1, 0)", row, col)
	}
}
