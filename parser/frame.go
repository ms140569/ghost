package parser

type Frame struct {
	command Cmd
	headers map[string]string
	data    []byte
}
