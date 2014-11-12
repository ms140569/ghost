package parser

type Cmd int

const (
	SEND Cmd = iota // 0
	SUBSCRIBE
	UNSUBSCRIBE
	BEGIN
	COMMIT
	ABORT
	ACK
	NACK
	DISCONNECT
	CONNECT // 9
	STOMP
	CONNECTED
	MESSAGE
	RECEIPT
	ERROR // 14

	// here are the synthetic commands no existing as real stomp frames
	// but used for certain purposes in the code

	COMMAND_NOT_RECOGNIZED // 15
	NOP                    // no operation (processor did not produce an answer
	HEARTBEAT              // send a heartbeat to the connection
)

var StompCommands = map[string]Cmd{
	"SEND":                   SEND,
	"SUBSCRIBE":              SUBSCRIBE,
	"UNSUBSCRIBE":            UNSUBSCRIBE,
	"BEGIN":                  BEGIN,
	"COMMIT":                 COMMIT,
	"ABORT":                  ABORT,
	"ACK":                    ACK,
	"NACK":                   NACK,
	"DISCONNECT":             DISCONNECT,
	"CONNECT":                CONNECT,
	"STOMP":                  STOMP,
	"CONNECTED":              CONNECTED,
	"MESSAGE":                MESSAGE,
	"RECEIPT":                RECEIPT,
	"ERROR":                  ERROR,
	"COMMAND_NOT_RECOGNIZED": COMMAND_NOT_RECOGNIZED,
	"NOP":       NOP,
	"HEARTBEAT": HEARTBEAT,
}

var requiredHeaders = map[Cmd][]string{
	CONNECT:     {"accept-version", "host"},
	STOMP:       {"accept-version", "host"},
	SEND:        {"destination"},
	SUBSCRIBE:   {"destination", "id"},
	UNSUBSCRIBE: {"id"},
	ACK:         {"id"},
	NACK:        {"id"},
	BEGIN:       {"transaction"},
	COMMIT:      {"transaction"},
	ABORT:       {"transaction"},
}

func (c Cmd) GetRequiredHeaders() []string {
	return requiredHeaders[c]
}

func CommandForString(operation string) Cmd {
	if cmd, ok := StompCommands[operation]; ok {
		return cmd
	}

	return COMMAND_NOT_RECOGNIZED
}

func (c Cmd) String() string {
	for key, value := range StompCommands {
		if value == c {
			return key
		}
	}

	return "Command not found"
}
