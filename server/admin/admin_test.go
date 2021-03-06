package admin

import (
	"testing"
)

func TestCommandParsingGood(t *testing.T) {

	// single

	verify(t, CommandScanner(wrap("status")))
	verify(t, CommandScanner(wrap("quit")))
	verify(t, CommandScanner(wrap("help")))

	// single + param

	verify(t, CommandScanner(wrap("show log")))

	// dest + sub

	verify(t, CommandScanner(wrap("dest list")))

	// dest + sub + param

	verify(t, CommandScanner(wrap("dest create gonzo")))
	verify(t, CommandScanner(wrap("dest delete gonzo")))
	verify(t, CommandScanner(wrap("dest stat gonzo")))

}

func TestCommandParsingBad(t *testing.T) {
	// NEGATIVE

	falsify(t, CommandScanner(wrap("no_valid_command"))) // garbage

	falsify(t, CommandScanner(wrap("quit param"))) // no params allowed

	falsify(t, CommandScanner(wrap("dest"))) // without subcommand

	falsify(t, CommandScanner(wrap("dest gonzo"))) // unknown subcommand

	falsify(t, CommandScanner(wrap("dest list something"))) // no params allowed

	falsify(t, CommandScanner(wrap("show"))) // params missing

}

func TestSubcommands(t *testing.T) {
	cmd := CommandScanner(wrap("dest list"))
	if cmd.sub != "list" {
		t.Fatalf("List destination command not parsed correctly, found sub: %s", cmd.sub)
	}

	cmd = CommandScanner(wrap("dest create frumpy"))
	if cmd.sub != "create" {
		t.Fatalf("Create destination command not parsed correctly, found sub: %s", cmd.sub)
	}

	if cmd.param != "frumpy" {
		t.Fatalf("Parameter not parsed correctly, found: %s", cmd.param)
	}

}

func wrap(data string) []byte {
	return []byte(data + "\n")
}

func verify(t *testing.T, token Shellcommand) {
	if token.name == UNDEF {
		t.Fatalf("Shellcommand was UNDEF")
	}
}

func falsify(t *testing.T, token Shellcommand) {
	if token.name != UNDEF {
		t.Fatalf("Shellcommand was *not* UNDEF but: %s", token.String())
	}

}
