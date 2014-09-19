package parser

import (
	"bytes"
	"github.com/ms140569/ghost/log"
	"strings"
)

const Separator string = ":"

type Frame struct {
	command Cmd
	headers map[string]string
	payload bytes.Buffer
}

func (f *Frame) addHeader(header string) {
	log.Debug("Adding header: |%s|", header)

	header = strings.TrimSuffix(header, "\r\n")
	header = strings.TrimSuffix(header, "\n")

	log.Debug("Adding header trimmed: |%s|", header)

	if f.headers == nil {
		log.Debug("Adding new header map")
		f.headers = make(map[string]string)
	}

	if strings.HasSuffix(header, Separator) {
		log.Debug("Adding header without value")
		f.headers[strings.TrimSuffix(header, Separator)] = ""
	} else {
		arr := strings.Split(header, Separator)
		log.Debug("Adding header with value")
		f.headers[arr[0]] = arr[1]
	}

}

func (f *Frame) dumpHeaders() {
	for k, v := range f.headers {
		log.Debug("HEADER, key: %s val: %s", k, v)
	}
}
