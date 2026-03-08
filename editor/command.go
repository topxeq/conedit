package editor

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
