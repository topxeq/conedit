package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/topxeq/conedit/editor"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Conedit - A lightweight command-line text editor")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  conedit                      Open editor for text input (default mode)")
		fmt.Println("  conedit [file.txt]           Open file for editing (immediate mode)")
		fmt.Println("  conedit -mode=file file.txt  Open file for editing (file mode)")
		fmt.Println("  conedit -mode=default        Open editor for text input")
		fmt.Println("  conedit -mode=immediate file.txt  Edit file with auto-save on exit")
		fmt.Println("  conedit -fromSSH -sshHost=... -sshUser=... -filePath=...  Edit remote file")
		fmt.Println()
		fmt.Println("Keyboard Shortcuts:")
		fmt.Println("  Ctrl+S  Save")
		fmt.Println("  Ctrl+K  Save As")
		fmt.Println("  Ctrl+X  Exit (immediate mode) / Confirm (default mode)")
		fmt.Println("  Ctrl+Q  Force Quit / Cancel (file mode)")
		fmt.Println("  Ctrl+W  Toggle Word Wrap")
		fmt.Println("  Ctrl+C  Copy")
		fmt.Println("  Ctrl+V  Paste")
		fmt.Println("  Ctrl+Z  Undo")
		fmt.Println("  Ctrl+Y  Redo")
		fmt.Println("  Ctrl+F  Find (supports regex)")
		fmt.Println("  Ctrl+H  Replace (supports regex)")
		fmt.Println("  Ctrl+G  Goto Line")
		fmt.Println()
		fmt.Println("Modes:")
		fmt.Println("  default    - Text input mode, no file operations (returns ok/cancel)")
		fmt.Println("  file       - Edit file, save/saveAs returns immediately (returns save/saveAs/cancel)")
		fmt.Println("  immediate  - Edit file, exit with Ctrl+X prompts to save (returns ok/cancel/error)")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -mode=MODE        Editor mode: default, file, immediate")
		fmt.Println("  -filePath=PATH    File path to edit")
		fmt.Println("  -fromSSH          Edit file on SSH server")
		fmt.Println("  -sshHost=HOST     SSH host")
		fmt.Println("  -sshPort=PORT     SSH port (default: 22)")
		fmt.Println("  -sshUser=USER     SSH username")
		fmt.Println("  -sshPass=PASS     SSH password")
		fmt.Println("  -sshKeyPath=PATH  SSH private key path")
		fmt.Println("  -mem              Force in-memory processing")
		fmt.Println("  -tmpPath=PATH     Custom temp directory for large files")
		os.Exit(0)
	}

	args := os.Args[1:]
	mode := ""
	filePath := ""

	// Parse explicit mode and file path
	for _, arg := range args {
		if strings.HasPrefix(arg, "-mode=") {
			mode = strings.TrimPrefix(arg, "-mode=")
		}
		if strings.HasPrefix(arg, "-filePath=") {
			filePath = strings.TrimPrefix(arg, "-filePath=")
		}
		if !strings.HasPrefix(arg, "-") && arg != "" && filePath == "" {
			filePath = arg
		}
	}

	// Auto-infer mode if not specified
	if mode == "" {
		if filePath != "" {
			mode = "immediate"
		} else {
			mode = "default"
		}
	}

	// Reconstruct args with explicit mode
	finalArgs := []string{"-mode=" + mode}
	if filePath != "" && !hasArg(args, "-filePath=") {
		finalArgs = append(finalArgs, "-filePath="+filePath)
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if arg == filePath && hasArg(args, "-filePath=") {
			continue
		}
		finalArgs = append(finalArgs, arg)
	}

	result := editor.ConsoleEditText("", finalArgs...)
	if result["error"] != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", result["error"])
		os.Exit(1)
	}

	status := result["status"]
	if status == "save" || status == "saveAs" || status == "ok" {
		if path, ok := result["path"].(string); ok && path != "" {
			fmt.Printf("File saved: %s\n", path)
		} else if result["text"] != "" {
			// immediate mode with changes saved
			fmt.Printf("Content saved\n")
		}
	}
}

func hasArg(args []string, prefix string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, prefix) {
			return true
		}
	}
	return false
}
