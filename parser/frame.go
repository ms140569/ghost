package parser

import (
	"log"
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

	header = strings.TrimSuffix(header, "\n")
	header = strings.TrimSuffix(header, "\r\n")

	log.Printf("Adding header trimmed: |%s|", header)

	if f.headers == nil {
		log.Println("Adding new header map")
		f.headers = make(map[string]string)
	}

	if strings.HasSuffix(header, Separator) {
		log.Println("Adding header without value")
		f.headers[strings.TrimSuffix(header, Separator)] = ""
	} else {
		arr := strings.Split(header, Separator)
		log.Println("Adding header with value")
		f.headers[arr[0]] = arr[1]
	}

}

func (f *Frame) dumpHeaders() {
	for k, v := range f.headers {
		log.Printf("HEADER, key: %s val: %s", k, v)
	}
}
