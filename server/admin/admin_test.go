package admin

import (
	"testing"
)

func TestCommandParsing(t *testing.T) {
	// CommandScanner(content [] byte) []Shellcommand

	input := "is nix"

	commandsParsed := CommandScanner([]byte(input))

	if len(commandsParsed) > 0 {
		t.Fatalf("This should produce *no* result: %s", input)
	}

	CommandScanner([]byte("status\n"))
	CommandScanner([]byte("dest\n"))
	CommandScanner([]byte("dest gonzo\n"))
	CommandScanner([]byte("dest list\n"))
	CommandScanner([]byte("show something\n"))
	CommandScanner([]byte("quit\n"))
	CommandScanner([]byte("quit param\n"))
}
