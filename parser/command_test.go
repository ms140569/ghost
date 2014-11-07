package parser

import (
	"testing"
)

func TestRequiredHeaders(t *testing.T) {
	if len(CONNECT.GetRequiredHeaders()) != 2 {
		t.Fatal("Wrong number of required headers for CONNECT")
	}
}
