package server

import (
	"github.com/ms140569/ghost/parser"
	"testing"
)

func TestUnknownFrame(t *testing.T) {
	frame := parser.NewFrame(parser.SEND)
	frame.AddHeader("crack:me")

	ProcessFrame(frame)

}
