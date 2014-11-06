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

type Parser struct {
	tokenIndex    int
	state         stateFn
	tokens        *[]Token
	frame         Frame
	startTime     time.Time
	err           error
	data          []byte
	bytesConsumed int
	nullIdx       int
}

func (p *Parser) next() Token {

	if p.tokenIndex >= len(*p.tokens) {
		log.Debug("EOF reached")
		return Token{name: EOF}
	}

	localTokens := *p.tokens
	tok := localTokens[p.tokenIndex]
	p.tokenIndex = p.tokenIndex + 1
	return tok
}

func (p *Parser) dumpTokens() {
	log.Debug("***********************************************")
	log.Debug("TOKENDUMP: Number of Tokens received: %d", len(*p.tokens))

	for number, token := range *p.tokens {

		var prefix string = ""

		switch token.name {
		case HEADER:
			prefix = "HEADER :"
		}

		log.Debug("%02d:%04d:%s%s\n", number, token.nextPos, prefix, token)
	}

	log.Debug("***********************************************")
}

func (p *Parser) nextPos() (int, error) {
	if p.tokenIndex > 0 {
		localTokens := *p.tokens
		tok := localTokens[(p.tokenIndex - 1)]
		return tok.nextPos, nil
	} else {
		return 0, errors.New("Could not get next position.")
	}

}

func (p *Parser) runMachine() {
	for state := startState; state != nil; {
		state = state(p)
	}
}

type stateFn func(*Parser) stateFn

func startState(p *Parser) stateFn {
	log.Debug("Start parsing bufffer, recording time.")
	p.startTime = time.Now()
	p.frame = Frame{}
	return getCommandState
}

func getCommandState(p *Parser) stateFn {
	log.Debug("In getCommandState")
	token := p.next()

	if token.name != COMMAND {
		p.err = errors.New("STOMP command not found.")
		return badExit
	}

	p.frame.command = token.value.(Cmd)

	token = p.next() // swallow EOL token to move on to the headers ...

	if token.name != EOL {
		p.err = errors.New("No EOL afer STOMP command found.")
		return badExit
	}

	return getHeadersState
}

func getHeadersState(p *Parser) stateFn {
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

func saveDataState(p *Parser) stateFn {
	log.Debug("saveDataState()")

	pos, err := p.nextPos()

	if err == nil {
		log.Debug("Data position: %d", pos)

		nullIdx := bytes.IndexByte(p.data[pos:], 0x00) // look for 0x00 terminator

		if nullIdx == -1 {
			p.err = errors.New("No null terminator found, bail out.")
			return badExit
		}

		log.Debug("Payload size: %d", nullIdx)

		p.frame.payload.Write(p.data[pos : pos+nullIdx])
		p.nullIdx = pos + nullIdx

		return swallowTrailingNewline
	} else {
		p.err = err
		return badExit
	}

}

func swallowTrailingNewline(p *Parser) stateFn {
	log.Debug("swallowTrailingNewline()")

	var pos int = p.nullIdx + 1 // ignore null byte itself.

	if pos == len(p.data) {
		p.bytesConsumed = pos
		return goodExit
	}

	for {
		b := p.data[pos]

		log.Debug("BYTE to swallow:%x on position: %d", b, pos)

		if b == 0x0a {
			log.Debug("Swallow a newline char.")
			pos = pos + 1
		} else {
			break
		}

		if pos == len(p.data) {
			break
		}

	}

	p.bytesConsumed = pos

	return goodExit
}

func badExit(p *Parser) stateFn {
	log.Error("badExit()")
	log.Error("Parsing error, last problem: %s", p.err)
	return cleanupAndExitMachine
}

func goodExit(p *Parser) stateFn {
	log.Debug("goodExit()")
	return cleanupAndExitMachine
}

func cleanupAndExitMachine(p *Parser) stateFn {
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

	var lastError error = nil
	var bytesRead int = 0

	for {
		number, frame, lastError := RunParser(data)

		log.Debug("Bytes read in this parsing run: %d", number)

		if lastError != nil {

			if globals.Config.Testmode {
				log.Fatal("Last parsing returned an error: %s", lastError.Error())
				os.Exit(1)
			}

			log.Debug("Last parsing returned an error: %s", lastError.Error())
			break
		}

		frames = append(frames, frame)
		bytesRead = bytesRead + number

		log.Debug("Number of Frames received: %d", len(frames))

		if len(data) <= 0 || len(data) <= bytesRead {
			break
		}

		data = data[bytesRead+1:]
	}

	return bytesRead, frames, lastError
}

func RunParser(data []byte) (int, Frame, error) {
	tokens := Scanner(data)

	if len(tokens) < 1 {
		msg := "Received no tokens, something is broken"

		if globals.Config.Testmode {
			log.Fatal("%s", msg)
			os.Exit(1)
		}

		log.Error(msg)
		return 0, Frame{}, errors.New(msg)
	}

	parser := CreateAndStartParser(data, tokens)

	parser.dumpTokens()

	return parser.bytesConsumed, parser.frame, parser.err
}

func CreateAndStartParser(data []byte, tokens []Token) Parser {
	parser := Parser{tokenIndex: 0, tokens: &tokens, data: data}
	parser.runMachine()
	return parser
}
