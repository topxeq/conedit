package editor

import "testing"

func TestCommand_Parse(t *testing.T) {
	tests := []struct {
		input    string
		expected Command
	}{
		{"\x13", CmdSave},       // Ctrl+S
		{"\x0b", CmdSaveAs},     // Ctrl+K
		{"\x03", CmdCopy},       // Ctrl+C
		{"\x16", CmdPaste},      // Ctrl+V
		{"\x1a", CmdUndo},       // Ctrl+Z
		{"\x19", CmdRedo},       // Ctrl+Y
		{"\x06", CmdFind},       // Ctrl+F
		{"\x08", CmdReplace},    // Ctrl+H
		{"\x07", CmdGotoLine},   // Ctrl+G
		{"\x18", CmdExit},       // Ctrl+X
		{"\x11", CmdForceQuit},  // Ctrl+Q
		{"\x17", CmdToggleWrap}, // Ctrl+W
	}
	for _, tt := range tests {
		result := ParseCommand(tt.input)
		if result != tt.expected {
			t.Errorf("ParseCommand(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
