package parser

import (
	"github.com/ms140569/ghost/log"
	"strings"
)

const Separator string = ":"

type Frame struct {
	command Cmd
	headers map[string]string
	data    []byte
}

func (f *Frame) addHeader(header string) {
	log.Printf("Adding header: |%s|", header)

	header = strings.TrimSuffix(header, "\r\n")
	header = strings.TrimSuffix(header, "\n")

	log.Printf("Adding header trimmed: |%s|", header)

	if f.headers == nil {
		log.Printf("Adding new header map")
		f.headers = make(map[string]string)
	}

	if strings.HasSuffix(header, Separator) {
		log.Printf("Adding header without value")
		f.headers[strings.TrimSuffix(header, Separator)] = ""
	} else {
		arr := strings.Split(header, Separator)
		log.Printf("Adding header with value")
		f.headers[arr[0]] = arr[1]
	}

}

func (f *Frame) dumpHeaders() {
	for k, v := range f.headers {
		log.Printf("HEADER, key: %s val: %s", k, v)
	}
}
