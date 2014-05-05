package main
import (
	"os"
	"os/exec"
	"testing"
	_"ostential"
)

func Test_ostential(t *testing.T) {
	cmd := exec.Command("go", "test", "ostential")
	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
}
