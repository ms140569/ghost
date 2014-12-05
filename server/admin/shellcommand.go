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

var helpMap = map[Name]string{
	STATUS: "status - Shows the system status.",
	HELP:   "help - Display help page.",
	QUIT:   "quit - Logout from system.",
	SHOW:   "show - Show various system configuration parameters.",
	DEST:   "dest - Destination related commands like list, create, delete, stat.",
}

func HelpForAllCommands() string {
	retVal := ""

	for _, msg := range helpMap {
		retVal = retVal + msg + "\n"
	}

	return retVal
}

func (t Shellcommand) Help() string {
	return helpMap[t.name]
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
