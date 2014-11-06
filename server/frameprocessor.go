package server

import (
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"net"
)

type LogicalConnection struct {
	isConnected bool
}

var connections map[net.Conn]LogicalConnection = make(map[net.Conn]LogicalConnection)

var frameQueue chan parser.Frame

func InitFrameQueue() {

	frameQueue = make(chan parser.Frame, 50)

	go FetchFrame()
}

func QueueFrame(conn net.Conn, frame parser.Frame) {
	log.Info("Queing frame: %s", frame.Command)

	frame.Connection = conn
	frameQueue <- frame
}

func FetchFrame() {
	for {
		frame := <-frameQueue
		log.Info("Processing single frame: %s", frame.Command.String())

		answer := ProcessFrame(frame)

		_, err := frame.Connection.Write([]byte(answer.Render()))

		if err != nil {
			frame.Connection.Close()
		}
	}
}

func ProcessFrame(input parser.Frame) parser.Frame {
	switch input.Command {
	case parser.CONNECT:
		return processConnect(input)
	default:
		return processDefault(input)
	}
}

func processConnect(frame parser.Frame) parser.Frame {
	log.Debug("processConnect")

	_, present := connections[frame.Connection]

	var answer parser.Frame

	if present {
		log.Debug("Connection know to server")
		answer = parser.NewFrame(parser.ERROR)
		answer.AddHeader("message:already connected.")
	} else {
		log.Debug("New connection, adding to map.")
		connections[frame.Connection] = LogicalConnection{isConnected: true}

		answer = parser.NewFrame(parser.CONNECTED)
		answer.AddHeader("not-used:value")
	}

	return answer
}

func processDefault(frame parser.Frame) parser.Frame {
	log.Debug("processDefault")
	answer := parser.NewFrame(parser.ACK)

	answer.AddHeader("not-used:value")
	answer.AddHeader("schmidtm:welcome back friend")

	return answer
}
