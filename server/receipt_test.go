package server

import (
	"github.com/ms140569/ghost/parser"
	"testing"
)

func TestSimpleReceipt(t *testing.T) {

	magic := "gonzo"

	frame := parser.NewFrame(parser.SEND)
	frame.AddHeader("destination:ok")
	frame.AddHeader("receipt:" + magic)

	answer := ProcessFrame(frame)

	if answer.Command != parser.RECEIPT {
		t.Fatalf("What have we got here? : %s", answer.Command)
	}

	if answer.GetHeader("receipt-id") != magic {
		t.Fatalf("Requested receipt but got none")
	}
}
