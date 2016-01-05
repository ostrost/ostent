package commands

import (
	"flag"
	"testing"

	"github.com/ostrost/ostent/commands/extpoints"
)

func Test(t *testing.T) {
	_ = NewWebserver(8080).AddCommandLine() // webserver :=
	Parse(flag.CommandLine, []string{"-b", ":8050"})
	bflag := flag.CommandLine.Lookup("b")
	if bflag == nil {
		t.Fatalf("Failed to find \"b\" flag.")
	}
	if bvalue, cmp := bflag.Value.String(), ":8050"; bvalue != cmp {
		t.Fatalf("Mismatch: %+v != %+v", bvalue, cmp)
	}
	errd, _ := ArgCommands()
	if errd {
		t.Fatalf("Must not continue after ArgCommands signaled to stop")
	}
	_, err := ParseCommand(
		[]extpoints.CommandHandler{},
		[]string{"nonexistentcommand"})
	if err == nil || err.Error() != "nonexistentcommand: No such command" {
		t.Fatalf("ParseCommand did not return expected error: %q", err)
	}
	clhs, err := ParseCommand(
		[]extpoints.CommandHandler{},
		[]string{"help"})
	if err != nil {
		t.Fatalf("ParseCommand failed: %q", err)
	}
	if v, cmp := len(clhs), 1; v != cmp {
		t.Fatalf("ParseCommand return unexpected (by length) result: %+v", clhs)
	}
}
