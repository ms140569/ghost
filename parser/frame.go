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
	log.Printf("Adding header: %s", header)
	if strings.HasSuffix(header, Separator) {
		f.headers[strings.TrimSuffix(header, Separator)] = ""
	} else {
		arr := strings.Split(header, Separator)
		f.headers[arr[0]] = arr[1]
	}

}
