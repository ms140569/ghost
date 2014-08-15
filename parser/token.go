package parser

type Name int

const (
	COMMAND Name = iota // value = Cmd - the command found
	HEADER              // value tuple: header = value - the header found with (optional) value
	DATA                // value []byte - the payload data as byte array
	NULL
	STRING
	OCTET
	EOL
	COLON
)

type Token struct {
	name  Name
	value interface{}
}

func (t Token) String() string {
	switch t.name {
	case COMMAND:
		return "COMMAND: " + (t.value.(Cmd)).String()
	case HEADER:
		return "HEADER: " + t.value.(string)
	case DATA:
		return "DATA, length:" + string(len(t.value.([]byte)))
	case NULL:
		return "NULL"
	case STRING:
		return "STRING: " + t.value.(string)
	case OCTET:
		return "OCTET"
	case EOL:
		return "EOL"
	case COLON:
		return "COLON"
	}

	return "Token not recognized"
}
