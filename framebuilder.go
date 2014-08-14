package main

import (
	"log"
)

type parser struct {
	pos          int
	state        stateFn
	tokens       *[]Token
	currentFrame *Frame
}

type stateFn func(*parser) stateFn

/*
func (p *parser) next() Token {
	return *(p.tokens)[p.pos]
} */

func (p *parser) run() {
	for state := startState; state != nil; {
		state = state(p)
	}
}

func startState(p *parser) stateFn {
	log.Println("In startState")
	return getCommandState
}

func getCommandState(p *parser) stateFn {
	log.Println("In getCommandState")
	return getHeadersState
}

func getHeadersState(p *parser) stateFn {
	log.Println("In getHeadersState")
	return nil
}

func ParseFrames(data []byte) []Frame {
	frames := []Frame{}

	tokens := Scanner(string(data))

	if len(tokens) < 1 {
		log.Println("Received no tokens, something is broken")
	}

	//for _, token := range tokens {
	//	log.Printf("%s\n", token)
	//}

	parser := parser{pos: 0, state: startState, tokens: &tokens}

	parser.run()

	return frames
}
