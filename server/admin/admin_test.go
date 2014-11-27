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
}
