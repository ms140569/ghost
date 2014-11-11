package server

import (
	"github.com/ms140569/ghost/parser"
	"testing"
)

func TestSimpleReceipt(t *testing.T) {
	frame := parser.NewFrame(parser.SEND)
	frame.AddHeader("destination:ok")
	frame.AddHeader("receipt:gonzo")

	answer := ProcessFrame(frame)

	if answer.Command != parser.RECEIPT {
		t.Fatalf("What have we got here? : %s", answer.Command)
	}

	if answer.GetHeader("receipt-id") != "gonzo" {
		t.Fatalf("Receipt request but got none")
	}
}
