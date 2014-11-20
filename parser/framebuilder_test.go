package parser

import (
	"testing"
)

func toFrameBody(s string) []byte {
	return []byte(s + "\x00")
}

func TestSimpleConnect(t *testing.T) {
	buffer := []byte("CONNECT\nsimple:header\n\n\x00")

	bytesConsumed, frames, err := ParseFrames(buffer, true)

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

	bytesConsumed, frames, err := ParseFrames(buffer, true)

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

func TestConnect1(t *testing.T) {
	buffer := toFrameBody(
		`CONNECT
simple:header
another:value
key:val
content-type:gonzo

`)

	bytesConsumed, frames, err := ParseFrames(buffer, true)

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

func parseAndVerify(t *testing.T, buffer []byte, size int, numberOfFrames int) {
	bytesConsumed, frames, err := ParseFrames(buffer, true)

	if bytesConsumed != size {
		t.Fatalf("Parser did not consume the correct number of bytes. Expected: %d, actually consumed: %d", size, bytesConsumed)
	}

	if len(frames) != numberOfFrames {
		t.Fatalf("Wrong number of Frames parsed. Expected: %d, actually parsed: %d", numberOfFrames, len(frames))
	}

	if err != nil {
		t.Fatalf("Parsing caused an error: %s", err)
	}

}

func TestConnect2(t *testing.T) {
	buffer := toFrameBody(
		"CONNECT\n" +
			"simple2:header2\n" +
			"string:plus\n" +
			"key:val\n" +
			"content-type:gonzo\n" +
			"\n")

	bytesConsumed, frames, err := ParseFrames(buffer, true)

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

func TestConnect3(t *testing.T) {
	buffer := toFrameBody(
		"CONNECT\n" +
			"simple2:header2\n" +
			"string:plus\n" +
			"key:val\n" +
			"content-type:gonzo\n" +
			"\n")

	parseAndVerify(t, buffer, 65, 1)
}

type parsingTests struct {
	data           string // input data to be parsed
	size           int    // number of bytes expected to be consumed
	numberOfFrames int    // number of frames to be parsed from input data
}

var connectTests = []parsingTests{
	{"CONNECT\nsimple:header\n\n\x00", 24, 1},
	{"CONNECT\nsimple2:header2\nkey:value\n\n\x00", 36, 1},
	{"CONNECT\nsimple2:header2\nkey:value\n\n\x00CONNECT\nsimple2:header2\nkey:value\n\n\x00", 72, 2},
}

func TestConnectsFromTable(t *testing.T) {

	for _, singleTest := range connectTests {
		parseAndVerify(t, []byte(singleTest.data), singleTest.size, singleTest.numberOfFrames)
	}

}
