package server

import (
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/parser"
	"github.com/ms140569/ghost/storage"
	"testing"
)

func TestMissingHeaderOnSend1(t *testing.T) {
	frame := parser.NewFrame(parser.SEND)
	frame.AddHeader("crack:me")

	answer := ProcessFrame(frame)

	if answer.Command != parser.ERROR {
		t.Fatal("Missing headers ought to produce error frames.")
	}
}

func TestGoodSend(t *testing.T) {

	globals.Config = globals.Configuration{}
	globals.Config.Storage = storage.Memory{}

	frame := parser.NewFrame(parser.SEND)
	frame.AddHeader("destination:ok")

	answer := ProcessFrame(frame)

	if answer.Command != parser.NOP {
		t.Fatalf("What have we got here? : %s", answer.Command)
	}
}

func TestMissingHeaderOnConnect(t *testing.T) {
	frame := parser.NewFrame(parser.CONNECT)
	frame.AddHeader("crack:me")

	answer := ProcessFrame(frame)

	if answer.Command != parser.ERROR {
		t.Fatal("Missing headers ought to produce error frames.")
	}
}

func TestSuccessfulConnect(t *testing.T) {
	frame := parser.NewFrame(parser.CONNECT)
	frame.AddHeader("accept-version:1.2")
	frame.AddHeader("host:localhost")

	answer := ProcessFrame(frame)

	if answer.Command != parser.CONNECTED {
		t.Fatalf("This should be a connected frame but was: %s", answer.Command)
	}

}
