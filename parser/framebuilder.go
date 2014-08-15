package parser

import (
	"log"
	"time"
	"errors"
)

// see: go/src/pkg/text/template/parse/lex.go
// or: http://rspace.googlecode.com/hg/slide/lex.html#landing-slide

type parser struct {
	pos    int
	state  stateFn
	tokens *[]Token
	frames *[]Frame
	startTime time.Time
	err error
}

type stateFn func(*parser) stateFn

func (p *parser) next() Token {
	localTokens := *p.tokens
	tok := localTokens[p.pos]
	p.pos = p.pos + 1
	return tok
}

func (p *parser) run() {
	for state := startState; state != nil; {
		state = state(p)
	}
}

func startState(p *parser) stateFn {
	log.Println("Start parsing bufffer, recording time.")
	p.startTime = time.Now()
	p.frames = &[]Frame{} // initialize empty array of frames ( frames will be appended here in getCommandState )
	return getCommandState
}

func getCommandState(p *parser) stateFn {
	log.Println("In getCommandState")
	token := p.next()

	if token.name != COMMAND {
		p.err = errors.New("STOMP command not found.")
		return badExit
	}

	*p.frames = append(*p.frames, Frame{command: token.value.(Cmd)})
 
	return getHeadersState
}

func getHeadersState(p *parser) stateFn {
	log.Println("In getHeadersState")
	log.Printf("Next Token would be: %s", p.next())
	return goodExit
}

func badExit(p *parser) stateFn {
	log.Println("badExit()")
	log.Printf("Parsing error, last problem: %s", p.err)
	dumpTokens(*p.tokens)
	return cleanupAndExitMachine
}

func goodExit(p *parser) stateFn {
	log.Println("goodExit()")
	return cleanupAndExitMachine
}

func cleanupAndExitMachine(p *parser) stateFn {
	log.Println("cleanupAndExitMachine()")
	log.Printf("Buffer parse-time: %v", time.Now().Sub(p.startTime))
	log.Printf("Number of Frames decoded: %d", len(*p.frames))
	return nil
}

func ParseFrames(data []byte) []Frame {
	frames := []Frame{}

	tokens := Scanner(string(data))

	if len(tokens) < 1 {
		log.Println("Received no tokens, something is broken")
	}

	parser := parser{pos: 0, state: startState, tokens: &tokens}

	parser.run()

	return frames
}

func dumpTokens(tokens []Token) {
	for _, token := range tokens {
		log.Printf("%s\n", token)
	}
}
