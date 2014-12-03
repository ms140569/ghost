package admin

import (
	"testing"
)

func TestCommandParsingGood(t *testing.T) {

	// single

	verify(CommandScanner(wrap("status")))
	verify(CommandScanner(wrap("quit")))
	verify(CommandScanner(wrap("help")))

	// single + param

	verify(CommandScanner(wrap("show log")))

	// dest + sub

	verify(CommandScanner(wrap("dest list")))

	// dest + sub + param

	verify(CommandScanner(wrap("dest create gonzo")))
	verify(CommandScanner(wrap("dest delete gonzo")))
	verify(CommandScanner(wrap("dest stat gonzo")))

}

func TestCommandParsingBad(t *testing.T) {
	// NEGATIVE

	falsify(CommandScanner(wrap("no_valid_command"))) // garbage

	falsify(CommandScanner(wrap("quit param"))) // no params allowed

	falsify(CommandScanner(wrap("dest"))) // without subcommand

	falsify(CommandScanner(wrap("dest gonzo"))) // unknown subcommand

	falsify(CommandScanner(wrap("dest list something"))) // no params allowed

	falsify(CommandScanner(wrap("show"))) // params missing

}

func wrap(data string) []byte {
	return []byte(data + "\n")
}

func verify(token Shellcommand) {
}

func falsify(token Shellcommand) {
}
