package parser

import (
	"testing"
)

func TestSimpleConnect(t *testing.T) {
	buffer := []byte("CONNECT\nsimple:header\n\n\x00")

	bytesConsumed, frames, err := ParseFrames(buffer)

	if bytesConsumed == 0 {
		t.Fatal("No bytes consumed")
	}

	if len(frames) == 0 {
		t.Fatal("No Frames returned")
	}

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Bytes consumed: %d, number of frames: %d", bytesConsumed, len(frames))
}

func TestMultilineStringLiteral(t *testing.T) {
	buffer := []byte(
		`CONNECT
simple:header
another:value

` + "\x00")

	bytesConsumed, frames, err := ParseFrames(buffer)

	if bytesConsumed == 0 {
		t.Fatal("No bytes consumed")
	}

	if len(frames) == 0 {
		t.Fatal("No Frames returned")
	}

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Bytes consumed: %d, number of frames: %d", bytesConsumed, len(frames))
}
