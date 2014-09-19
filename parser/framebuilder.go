package parser

import (
	"bytes"
	"errors"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"os"
	"time"
)

// see: go/src/pkg/text/template/parse/lex.go
// or: http://rspace.googlecode.com/hg/slide/lex.html#landing-slide

type parser struct {
	pos       int
	state     stateFn
	tokens    *[]Token
	frames    *[]Frame
	startTime time.Time
	err       error
	data      []byte
}

type stateFn func(*parser) stateFn

func (p *parser) getLastFrame() (*Frame, error) {
	numberOfFrames := len(*p.frames)

	log.Debug("Number of frames: %d", numberOfFrames)

	if numberOfFrames < 1 {
		return &Frame{}, errors.New("There is no last frame.")
	} else {
		log.Debug("Fetching frame...")
		lFrames := *p.frames
		return &lFrames[numberOfFrames-1], nil
	}
}

func (p *parser) next() Token {

	if p.pos >= len(*p.tokens) {
		log.Debug("EOF reached")
		return Token{name: EOF}
	}

	localTokens := *p.tokens
	tok := localTokens[p.pos]
	p.pos = p.pos + 1
	return tok
}

func (p *parser) nextPos() (int, error) {
	if p.pos > 0 {
		localTokens := *p.tokens
		tok := localTokens[(p.pos - 1)]
		return tok.nextPos, nil
	} else {
		return 0, errors.New("Could not get next position.")
	}

}

func (p *parser) run() {
	for state := startState; state != nil; {
		state = state(p)
	}
}

func startState(p *parser) stateFn {
	log.Debug("Start parsing bufffer, recording time.")
	p.startTime = time.Now()
	p.frames = &[]Frame{} // initialize empty array of frames ( frames will be appended here in getCommandState )
	return getCommandState
}

func getCommandState(p *parser) stateFn {
	log.Debug("In getCommandState")
	token := p.next()

	if token.name != COMMAND {
		p.err = errors.New("STOMP command not found.")
		return badExit
	}

	*p.frames = append(*p.frames, Frame{command: token.value.(Cmd)})

	// swallow EOL token to move on to the headers ...
	// FIXME: this ought to be done in the grammar file
	token = p.next()

	if token.name != EOL {
		p.err = errors.New("No EOL afer STOMP command found.")
		return badExit
	}

	return getHeadersState
}

func getHeadersState(p *parser) stateFn {
	log.Debug("In getHeadersState")

	token := p.next()

	if token.name == EOL {
		log.Debug("Empty header set found. Moving on to data section.")
		return saveDataState
	}

	for token.name == HEADER {
		log.Debug("Header found: %s", token)
		lastFrame, err := p.getLastFrame()

		if err == nil {
			lastFrame.addHeader(token.String())
		} else {
			log.Error("Could not find last Frame.")
			return badExit
		}

		token = p.next()
		if token.name == EOL {
			return saveDataState
		}
	}

	p.err = errors.New("Headers corrupt.")
	return badExit
}

func saveDataState(p *parser) stateFn {
	log.Debug("saveDataState()")

	pos, err := p.nextPos()

	if err == nil {
		log.Debug("Data position: %d", pos)

		slice := p.data[pos:]
		nullIdx := bytes.IndexByte(slice, 0x00) // look for 0x00 terminator

		if nullIdx == -1 {
			p.err = errors.New("No null terminator found, bail out.")
			return badExit
		}

		log.Debug("NUL byte position: %d", nullIdx)

		lastFrame, err := p.getLastFrame()

		if err == nil {
			lastFrame.payload.Write(p.data[pos : pos+nullIdx])
		} else {
			log.Error("Could not find last Frame.")
			return badExit
		}

		return goodExit
	} else {
		p.err = err
		return badExit
	}

}

func badExit(p *parser) stateFn {
	log.Error("badExit()")
	log.Error("Parsing error, last problem: %s", p.err)
	dumpTokens(*p.tokens)
	return cleanupAndExitMachine
}

func goodExit(p *parser) stateFn {
	log.Debug("goodExit()")
	dumpTokens(*p.tokens)
	return cleanupAndExitMachine
}

func cleanupAndExitMachine(p *parser) stateFn {
	log.Debug("cleanupAndExitMachine()")
	log.Debug("Buffer parse-time: %v", time.Now().Sub(p.startTime))
	log.Debug("Number of Frames decoded: %d", len(*p.frames))

	lastFrame, _ := p.getLastFrame()

	log.Debug("Number of headers of last Frame: %d", len(lastFrame.headers))
	log.Debug("Number of payload: %d", lastFrame.payload.Len())

	lastFrame.dumpHeaders()

	if globals.Config.Testmode {
		log.Info("Running in testmode, exit.")
		if p.err == nil {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	return nil
}

func ParseFrames(data []byte) []Frame {

	tokens := Scanner(data)

	if len(tokens) < 1 {
		log.Error("Received no tokens, something is broken")
	}

	parser := parser{pos: 0, state: startState, tokens: &tokens, data: data}

	parser.run()

	return nil
}

func dumpTokens(tokens []Token) {
	log.Debug("***********************************************")
	for number, token := range tokens {

		var prefix string = ""

		switch token.name {
		case HEADER:
			prefix = "HEADER :"
		}

		log.Debug("%02d:%04d:%s%s\n", number, token.nextPos, prefix, token)
	}

	log.Debug("***********************************************")
}
