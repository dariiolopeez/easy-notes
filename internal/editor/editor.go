package editor

import (
	"os"
	"os/exec"
)

func Open(editorCmd, filePath string) error {
	cmd := exec.Command("sh", "-c", editorCmd+` "$0"`, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
