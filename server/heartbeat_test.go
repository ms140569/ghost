package server

import (
	"testing"
)

func TestHeartbeatParsing(t *testing.T) {

	var out int
	var in int
	var err error

	// FAILURE *****************************************

	out, in, err = parseHeartbeat("")

	if out != -1 || in != -1 || err == nil {
		t.Fatalf("This ought to break")
	}

	out, in, err = parseHeartbeat("nix")

	if out != -1 || in != -1 || err == nil {
		t.Fatalf("This ought to break")
	}

	out, in, err = parseHeartbeat("x,y")

	if out != -1 || in != -1 || err == nil {
		t.Fatalf("This ought to break")
	}

	out, in, err = parseHeartbeat("-1,-2")

	if out != -1 || in != -1 || err == nil {
		t.Fatalf("This ought to break")
	}

	out, in, err = parseHeartbeat(",,,,,")

	if out != -1 || in != -1 || err == nil {
		t.Fatalf("This ought to break")
	}

	// SUCESS **********************************************

	out, in, err = parseHeartbeat("0,0")

	if out == -1 || in == -1 || err != nil {
		t.Fatalf("This should be OK")
	}

	out, in, err = parseHeartbeat("1,2")

	if out == -1 || in == -1 || err != nil {
		t.Fatalf("This should be OK")
	}

	out, in, err = parseHeartbeat("22,3")

	if out == -1 || in == -1 || err != nil {
		t.Fatalf("This should be OK")
	}

	out, in, err = parseHeartbeat("99999,8888888")

	if out == -1 || in == -1 || err != nil {
		t.Fatalf("This should be OK")
	}

}
