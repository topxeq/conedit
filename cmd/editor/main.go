package main

import (
	"fmt"
	"os"

	"github.com/topxeq/conedit/editor"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Console Editor - A lightweight command-line text editor")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  console_editor [file.txt]           Open file for editing")
		fmt.Println("  console_editor -filePath=file.txt   Open file for editing")
		fmt.Println("  console_editor -fromSSH -sshHost=... -sshUser=... -filePath=...  Edit remote file")
		fmt.Println()
		fmt.Println("Keyboard Shortcuts:")
		fmt.Println("  Ctrl+S  Save")
		fmt.Println("  Ctrl+K  Save As")
		fmt.Println("  Ctrl+X  Exit")
		fmt.Println("  Ctrl+Q  Force Quit")
		fmt.Println("  Ctrl+W  Toggle Word Wrap")
		fmt.Println("  Ctrl+C  Copy")
		fmt.Println("  Ctrl+V  Paste")
		fmt.Println("  Ctrl+Z  Undo")
		fmt.Println("  Ctrl+Y  Redo")
		fmt.Println("  Ctrl+F  Find (supports regex)")
		fmt.Println("  Ctrl+H  Replace (supports regex)")
		fmt.Println("  Ctrl+G  Goto Line")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -filePath=PATH      File path to edit")
		fmt.Println("  -fromSSH            Edit file on SSH server")
		fmt.Println("  -sshHost=HOST       SSH host")
		fmt.Println("  -sshPort=PORT       SSH port (default: 22)")
		fmt.Println("  -sshUser=USER       SSH username")
		fmt.Println("  -sshPass=PASS       SSH password")
		fmt.Println("  -sshKeyPath=PATH    SSH private key path")
		fmt.Println("  -mem                Force in-memory processing")
		fmt.Println("  -tmpPath=PATH       Custom temp directory for large files")
		os.Exit(0)
	}

	result := editor.ConsoleEditText("", os.Args[1:]...)
	if result["error"] != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", result["error"])
		os.Exit(1)
	}

	status := result["status"]
	if status == "save" || status == "saveAs" {
		if path, ok := result["path"].(string); ok {
			fmt.Printf("File saved: %s\n", path)
		}
	}
}
