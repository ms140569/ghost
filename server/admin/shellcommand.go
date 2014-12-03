package admin

import "strings"

type Name int

const (
	STATUS Name = iota
	HELP
	QUIT
	DEST
	SHOW
	UNDEF
)

type Shellcommand struct {
	name  Name
	sub   string
	param string
}

func ShellCommandNameForString(input string) Name {
	switch strings.ToLower(input) {
	case "status":
		return STATUS
	case "help":
		return HELP
	case "quit":
		return QUIT
	case "dest":
		return DEST
	case "show":
		return SHOW
	default:
		return UNDEF

	}
}

func (t Shellcommand) String() string {

	switch t.name {
	case STATUS:
		return "STATUS"
	case HELP:
		return "HELP"
	case QUIT:
		return "QUIT"
	case SHOW:
		return "SHOW"
	case DEST:
		return "DEST"
	}

	return "Command not recognized"
}
