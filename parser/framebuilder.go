package parser

import (
	"log"
)

// see: go/src/pkg/text/template/parse/lex.go
// or: http://rspace.googlecode.com/hg/slide/lex.html#landing-slide

type parser struct {
	pos    int
	state  stateFn
	tokens *[]Token
	frames *[]Frame
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
	log.Println("In startState")
	log.Printf("Next Token would be: %s", p.next())
	return getCommandState
}

func getCommandState(p *parser) stateFn {
	log.Println("In getCommandState")
	log.Printf("Next Token would be: %s", p.next())
	return getHeadersState
}

func getHeadersState(p *parser) stateFn {
	log.Println("In getHeadersState")
	log.Printf("Next Token would be: %s", p.next())
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
