package parser

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"net"
	"strings"
)

const Separator string = ":"

type Frame struct {
	Command    Cmd
	headers    map[string]string
	payload    bytes.Buffer
	Connection net.Conn
}

func NewFrame(cmd Cmd) Frame {
	frame := Frame{Command: cmd}
	frame.headers = make(map[string]string)
	return frame
}

func (f *Frame) AddHeader(header string) error {
	header = strings.TrimSuffix(header, "\r\n")
	header = strings.TrimSuffix(header, "\n")

	if f.headers == nil {
		log.Debug("Adding new header map")
		f.headers = make(map[string]string)
	}

	// enforcing header related limitations.

	if len(f.headers) >= globals.MaxHeaderSize {
		return errors.New("Maximum number of headers reached")
	}

	if strings.HasSuffix(header, Separator) {
		key := strings.TrimSuffix(header, Separator)

		if len(key) > globals.MaxHaederKeyLength {
			return errors.New("Header key too long.")
		}

		f.headers[key] = ""
	} else {
		arr := strings.Split(header, Separator)
		key := arr[0]
		val := arr[1]

		if len(key) > globals.MaxHaederKeyLength {
			return errors.New("Header key too long.")
		}

		if len(val) > globals.MaxHaederValLength {
			return errors.New("Header value too long.")
		}

		f.headers[key] = val
	}

	return nil
}

func (f *Frame) dumpHeaders() {
	for k, v := range f.headers {
		log.Debug("HEADER, key: %s val: %s", k, v)
	}
}

func (f *Frame) Render() string {
	// Command
	retVal := fmt.Sprintf(f.Command.String()) + "\n"
	// headers

	for k, v := range f.headers {
		retVal = retVal + fmt.Sprintf("%s:%s\n", k, v)
	}
	// data

	if f.payload.Len() > 0 {
		retVal = retVal + fmt.Sprintf("%s", f.payload.String())
	}

	// NULL
	retVal = retVal + fmt.Sprintf("\x00\n")

	return retVal
}
