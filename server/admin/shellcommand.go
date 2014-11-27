package admin

type Name int

const (
	STATUS Name = iota
	HELP
	QUIT
	DEST
	EOL
)

type Shellcommand struct {
	name  Name
	sub   string
	value interface{}
}

func (t Shellcommand) String() string {

	switch t.name {
	case STATUS:
		return "STATUS"
	case HELP:
		return "HELP"
	case QUIT:
		return "QUIT"
	case DEST:
		return "DEST"
	}

	return "Command not recognized"
}
