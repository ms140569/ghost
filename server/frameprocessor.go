package server

import (
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"net"
)

var frameQueue chan parser.Frame
var conn net.Conn

func InitFrameQueue() {

	frameQueue = make(chan parser.Frame)

	go ProcessFrame()
}

func QueueFrame(conn net.Conn, frame parser.Frame) {
	log.Info("Queing frame.")

	frame.Connection = conn
	frameQueue <- frame
}

func ProcessFrame() {
	frame := <-frameQueue
	log.Info("Processing single frame.")

	answer := parser.NewFrame(parser.CONNECTED)
	answer.AddHeader("not-used:value")

	_, err := frame.Connection.Write([]byte(answer.Render()))

	if err != nil {
		frame.Connection.Close()
	}

}
