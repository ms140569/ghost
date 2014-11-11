package parser

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"net"
	"strconv"
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

func (f *Frame) HasHeader(header string) bool {
	_, present := f.headers[header]
	return present
}

func (f *Frame) GetHeader(header string) string {
	return f.headers[header]
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

	// dont forget trailing newline to separate headers.
	retVal = retVal + fmt.Sprintf("\n")

	// data

	if f.payload.Len() > 0 {
		retVal = retVal + fmt.Sprintf("%s", f.payload.String())
	}

	// NULL
	retVal = retVal + fmt.Sprintf("\x00\n")

	return retVal
}

/*
   Return the content length possibly given by the "content-length" header.
   Options are:

   - header NOT given -> -1
   - header given but empty -> error
   - header given but not parsable ( negative, ascii, etc.) -> error
   - header given with accurate positive integer -> integer
*/

func (f *Frame) GetContentLength() (int, error) {
	if f.HasHeader("content-length") {

		contentLengthAsString := f.GetHeader("content-length")

		if len(contentLengthAsString) > 0 {
			contentLength, err := strconv.Atoi(contentLengthAsString)
			if err != nil {
				return -2, err
			} else {
				if contentLength < 0 {
					return -2, errors.New("negative content-length value given.")
				}
				return contentLength, nil
			}
		} else {
			return -2, errors.New("content-length header given, but empty")
		}
	} else {
		return -1, nil
	}

}
