package editor

import "testing"

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
