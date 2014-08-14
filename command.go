package main

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
	ERROR
	COMMAND_NOT_RECOGNIZED // 15
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
