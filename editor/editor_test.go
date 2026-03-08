package editor

import (
	"os"
	"testing"
)

func TestModeBehavior_DefaultMode_ReturnsOkOnSave(t *testing.T) {
	t.Skip("Interactive test - requires manual verification")

	opts := []string{"-mode=default"}
	result := ConsoleEditText("test content", opts...)

	if result["status"] != "cancel" {
		t.Errorf("default mode with no interaction should return cancel, got %v", result["status"])
	}
}

func TestModeBehavior_FileMode_RequiresFilePath(t *testing.T) {
	t.Skip("Interactive test - requires manual verification")

	opts := []string{"-mode=file"}
	result := ConsoleEditText("", opts...)

	if result["status"] != "cancel" {
		t.Errorf("file mode with no file should return cancel, got %v", result["status"])
	}
}

func TestModeBehavior_ImmediateMode_RequiresFilePath(t *testing.T) {
	t.Skip("Interactive test - requires manual verification")

	opts := []string{"-mode=immediate"}
	result := ConsoleEditText("", opts...)

	if result["status"] != "cancel" {
		t.Errorf("immediate mode with no file should return cancel, got %v", result["status"])
	}
}

func TestParseOpts_Mode(t *testing.T) {
	opts := ParseOpts([]string{"-mode=default", "-filePath=/tmp/test.txt"})

	if opts["mode"] != "default" {
		t.Errorf("expected mode=default, got %v", opts["mode"])
	}
	if opts["filePath"] != "/tmp/test.txt" {
		t.Errorf("expected filePath=/tmp/test.txt, got %v", opts["filePath"])
	}
}

func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("TERM", "xterm-256color")

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}
