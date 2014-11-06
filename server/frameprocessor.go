package server

import (
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"net"
)

var frameQueue chan parser.Frame
var conn net.Conn

func InitFrameQueue() {

	frameQueue = make(chan parser.Frame, 50)

	go ProcessFrame()
}

func QueueFrame(conn net.Conn, frame parser.Frame) {
	log.Info("Queing frame: %s", frame.Command)

	frame.Connection = conn
	frameQueue <- frame
}

func ProcessFrame() {
	for {
		frame := <-frameQueue
		log.Info("Processing single frame: %s", frame.Command.String())

		answer := parser.NewFrame(parser.CONNECTED)
		answer.AddHeader("not-used:value")

		_, err := frame.Connection.Write([]byte(answer.Render()))

		if err != nil {
			frame.Connection.Close()
		}
	}
}
