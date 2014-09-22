package parser

import (
	"bytes"
	"errors"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"strings"
)

const Separator string = ":"

type Frame struct {
	command Cmd
	headers map[string]string
	payload bytes.Buffer
}

func (f *Frame) addHeader(header string) error {
	header = strings.TrimSuffix(header, "\r\n")
	header = strings.TrimSuffix(header, "\n")

	log.Debug("Adding header: |%s|", header)

	if f.headers == nil {
		log.Debug("Adding new header map")
		f.headers = make(map[string]string)
	}

	// enforcing header related limitations.

	if len(f.headers) >= globals.MaxHeaderSize {
		return errors.New("Maximum number of headers reached")
	}

	if strings.HasSuffix(header, Separator) {
		log.Debug("Adding header without value")

		key := strings.TrimSuffix(header, Separator)

		if len(key) > globals.MaxHaederKeyLength {
			return errors.New("Header key too long.")
		}

		f.headers[key] = ""
	} else {
		arr := strings.Split(header, Separator)
		key := arr[0]
		val := arr[1]

		log.Debug("Adding header with value")

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
