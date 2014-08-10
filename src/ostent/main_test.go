package main
import (
	"os"
	"os/exec"
	"testing"
	_"ostent"
)

func Test_ostent(t *testing.T) {
	cmd := exec.Command("go", "test", ".") // "ostent"
	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}
}
