package server

import (
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"github.com/twinj/uuid"
	"net"
	"os"
)

type LogicalConnection struct {
	isConnected  bool
	isProtocol12 bool
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

		if answer.Command != parser.COMMAND_NOT_RECOGNIZED {
			_, err := frame.Connection.Write([]byte(answer.Render()))

			if err != nil {
				frame.Connection.Close()
			}
		}

		// if the answer is an error frame, close the connection
		if answer.Command == parser.ERROR {
			frame.Connection.Close()
		}

		// do we have to send a receipt after the Frame?

		if frame.HasHeader("receipt") {
			receipt := parser.NewFrame(parser.RECEIPT)
			receipt.AddHeader("receipt-id:" + frame.GetHeader("receipt"))

			_, err := frame.Connection.Write([]byte(receipt.Render()))

			if err != nil {
				frame.Connection.Close()
			}
		}
	}
}

func ProcessFrame(frame parser.Frame) parser.Frame {

	// check for required frame headers

	for _, header := range frame.Command.GetRequiredHeaders() {
		if !frame.HasHeader(header) {
			msg := "Missing header - " + header

			if globals.Config.Testmode {
				log.Fatal("%s", msg)
				os.Exit(1)
			}

			log.Error(msg)

			// produce error Frame
			return createErrorFrameWithMessage(msg)
		}
	}

	// Check that receipts are only present on any client frame but CONNECT

	if frame.HasHeader("receipt") {
		if frame.Command == parser.CONNECT {
			return createErrorFrameWithMessage("CONNECT frames must not contain receipt headers.")
		}
	}

	// dispatch frame to handler function

	switch frame.Command {
	case parser.CONNECT:
		return processConnect(frame)
	case parser.SEND:
		return processSend(frame)
	default:
		return processDefault(frame)
	}
}

func createErrorFrameWithMessage(msg string) parser.Frame {
	answer := parser.NewFrame(parser.ERROR)
	answer.AddHeader("message:" + msg)
	return answer
}

func processConnect(frame parser.Frame) parser.Frame {
	log.Debug("processConnect")

	_, present := connections[frame.Connection]

	var answer parser.Frame

	if present {
		log.Debug("Connection know to server")
		answer = createErrorFrameWithMessage("already connected")
	} else {
		log.Debug("New connection, adding to map.")
		connections[frame.Connection] = LogicalConnection{isConnected: true}

		answer = parser.NewFrame(parser.CONNECTED)
		answer.AddHeader("version:1.2")

		// generate a session id
		sessionId := uuid.NewV4()
		answer.AddHeader("session:" + sessionId.String())

		// server versioin
		answer.AddHeader("server:" + globals.GetServerVersionString())
	}

	return answer
}

func processSend(frame parser.Frame) parser.Frame {
	log.Debug("processSend")
	return parser.NewFrame(parser.COMMAND_NOT_RECOGNIZED)
}

func processDefault(frame parser.Frame) parser.Frame {
	log.Debug("processDefault")
	return createErrorFrameWithMessage("Unknown Frame")
}
