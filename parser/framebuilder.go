package parser

import (
	"bytes"
	"errors"
	"github.com/ms140569/ghost/log"
	"time"
)

// see: go/src/pkg/text/template/parse/lex.go
// or: http://rspace.googlecode.com/hg/slide/lex.html#landing-slide

type parser struct {
	pos       int
	state     stateFn
	tokens    *[]Token
	frame     Frame
	startTime time.Time
	err       error
	data      []byte
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

type stateFn func(*parser) stateFn

func startState(p *parser) stateFn {
	log.Debug("Start parsing bufffer, recording time.")
	p.startTime = time.Now()
	p.frame = Frame{} 
	return getCommandState
}

func getCommandState(p *parser) stateFn {
	log.Debug("In getCommandState")
	token := p.next()

	if token.name != COMMAND {
		p.err = errors.New("STOMP command not found.")
		return badExit
	}

	//*p.frames = append(*p.frames, Frame{command: token.value.(Cmd)})

	p.frame.command = token.value.(Cmd)

	// swallow EOL token to move on to the headers ...
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

			err := p.frame.AddHeader(token.String())

			if err != nil {
				p.err = err
				log.Error(p.err.Error())
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

		nullIdx := bytes.IndexByte(p.data[pos:], 0x00) // look for 0x00 terminator

		if nullIdx == -1 {
			p.err = errors.New("No null terminator found, bail out.")
			return badExit
		}

		log.Debug("NUL byte position: %d", nullIdx)

		p.frame.payload.Write(p.data[pos : pos+nullIdx])

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

	log.Debug("Number of headers found: %d", len(p.frame.headers))
	log.Debug("Number of payload: %d", p.frame.payload.Len())

	p.frame.dumpHeaders()

	return nil
}
/*
Parses the input data given in the slice into a slice of Frames.
*/
func ParseFrames(data []byte) (int, []Frame, error) {

	frames := []Frame{}

	for {
		number, frame, lastError := RunParser(data)

		log.Debug("Bytes read: %d", number)
		
		if lastError != nil {
			log.Debug("Last parsing returned an error: %s", lastError.Error())
		} 

		// http://stackoverflow.com/questions/20240179/nil-detection-in-golang
		frames = append(frames, frame)		

		_ = frame
		_ = lastError


		break
	}

	log.Debug("Number of Frames received: %d", len(frames))

	return 0, nil, nil
}

func RunParser(data []byte) (int, Frame, error) {
	tokens := Scanner(data)

	if len(tokens) < 1 {
		msg := "Received no tokens, something is broken"
		log.Error(msg)
		return 0, Frame{}, errors.New(msg)
	}

	parser := parser{pos: 0, tokens: &tokens, data: data}

	parser.run()

	return 0, parser.frame, parser.err
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
