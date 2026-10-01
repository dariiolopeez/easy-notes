package editor

import (
	"os"
	"os/exec"
	"strings"
)

func Open(editorCmd, filePath string) error {
	parts := strings.Fields(editorCmd)
	args := append(parts[1:], filePath)
	cmd := exec.Command(parts[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
