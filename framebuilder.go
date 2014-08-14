package main

import (
	"fmt"
	//"log"
)

func ParseFrames(data []byte) []Frame {
	frames := []Frame{}

	tokens := Scanner(string(data))

	if len(tokens) < 1 {
		fmt.Println("Received no tokens, something is broken")
	}

	for _, token := range tokens {
		fmt.Printf("%s\n", token)
	}

	return frames
}
