package server

import (
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"net"
	"os"
)

var inboundFrameQueue chan parser.Frame
var outboundFrameQueue chan parser.Frame

func InitFrameQueues() {

	inboundFrameQueue = make(chan parser.Frame, globals.QueueSizeInbound)
	outboundFrameQueue = make(chan parser.Frame, globals.QueueSizeOutbound)

	if !globals.Config.Provider.Initialize() {
		log.Fatal("Unable to initialize the storage")
		os.Exit(3)
	}

	go FetchFrame()
	go FrameWriter()

	// Support for heartbeating
	go HeartBeatChecker()
}

func QueueFrame(conn net.Conn, frame parser.Frame) {
	log.Info("Queing frame: %s", frame.Command)

	frame.Connection = conn
	inboundFrameQueue <- frame
}

func FetchFrame() {
	for {
		frame := <-inboundFrameQueue
		log.Info("Processing single frame: %s", frame.Command.String())

		answer := ProcessFrame(frame)

		if answer.Command != parser.NOP {
			writeAnswer(frame.Connection, answer)
		}

		// if the answer is an error frame, close the connection
		if answer.Command == parser.ERROR {
			frame.Connection.Close()
		}
	}
}

func FrameWriter() {
	for {
		frame := <-outboundFrameQueue

		log.Debug("FrameWriter sending Frame: %s", frame.Command.String())

		// we either send a full-fledged frame or a heartbeat

		var output []byte

		if frame.Command == parser.HEARTBEAT {
			output = []byte("\x0A")
		} else {
			output = []byte(frame.Render())
		}

		if frame.Connection == nil {
			log.Error("No connection to write to in FrameWriter")
		} else {
			_, err := frame.Connection.Write(output)

			if err != nil || frame.Command == parser.ERROR {
				removeSessionFromMaps(frame.Connection)
				frame.Connection.Close()
			}
		}

	}
}

/*
Here we write the answer frame back to the connection. A couple of strategies will evolve over time:

- simple send to the connection ( not save )
- send to the connection using some lock sync.Lock/Unlock ( better )
- passing output to a channel to send it
*/

func writeAnswer(conn net.Conn, answer parser.Frame) {

	answer.Connection = conn
	outboundFrameQueue <- answer

	/*
		_, err := conn.Write([]byte(answer.Render()))

		if err != nil {
			conn.Close()
		}
	*/
}

func ProcessFrame(frame parser.Frame) parser.Frame {

	// We do the bookkeepting for heartbeating first, knowing that this means
	// that bounced or invalid frames keep connections alive.
	// I might change my mind on this in the long run.

	updateKeepaliveRecords(frame.Connection)

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
		processSend(frame)
	case parser.SUBSCRIBE:
		processSubscribe(frame)
	case parser.UNSUBSCRIBE:
		processUnsubscribe(frame)
	case parser.HEARTBEAT:
		processHeartbeat(frame)
	case parser.DISCONNECT:
		processDisconnect(frame)
	default:
		return processDefault(frame)
	}

	// do we have to send a receipt after the Frame?

	if frame.HasHeader("receipt") {
		receipt := parser.NewFrame(parser.RECEIPT)
		receipt.AddHeader("receipt-id:" + frame.GetHeader("receipt"))

		return receipt
	}

	// sending an empty, dummy frame to indicate that there is no reply.
	return parser.NewFrame(parser.NOP)
}

func createErrorFrameWithMessage(msg string) parser.Frame {
	answer := parser.NewFrame(parser.ERROR)
	answer.AddHeader("message:" + msg)

	log.Error("Created ERROR frame with message: " + msg)
	return answer
}

func processConnect(frame parser.Frame) parser.Frame {
	log.Debug("processConnect")

	_, present := sessions[frame.Connection]

	var answer parser.Frame

	if present {
		log.Debug("Session know to server")
		answer = createErrorFrameWithMessage("already connected")
	} else {

		// FIXME: here should go the code to check for login and passcode

		if frame.HasHeader("login") && frame.HasHeader("passcode") {
			msg := "login and passcode not supported yet."
			return createErrorFrameWithMessage(msg)
		}

		// create a new session and store it in the session map for further reuse

		log.Debug("New session, adding to map.")

		session := NewSession(frame.Connection)

		sessions[frame.Connection] = &session

		answer = parser.NewFrame(parser.CONNECTED)
		answer.AddHeader("version:1.2")

		// Analyze and setup heartbeating ...

		if frame.HasHeader("heart-beat") {
			heartbeatAnswer, err := initializeHeartbeatingForConnection(frame, &session)

			if err != nil {
				return createErrorFrameWithMessage(err.Error())
			} else {
				answer.AddHeader(heartbeatAnswer)
			}
		}

		answer.AddHeader("session:" + session.id)

		// server versioin
		answer.AddHeader("server:" + globals.GetServerVersionString())
	}

	return answer
}

func processSend(frame parser.Frame) parser.Frame {
	log.Debug("processSend")
	return parser.NewFrame(parser.NOP)
}

func processSubscribe(frame parser.Frame) parser.Frame {
	log.Debug("processSubscribe")
	return parser.NewFrame(parser.NOP)
}

func processUnsubscribe(frame parser.Frame) parser.Frame {
	log.Debug("processUnsubscribe")
	return parser.NewFrame(parser.NOP)
}

func processDisconnect(frame parser.Frame) parser.Frame {
	log.Debug("processDisconnect")
	return parser.NewFrame(parser.NOP)
}

func processHeartbeat(frame parser.Frame) parser.Frame {
	log.Debug("processHeartbeat")
	return parser.NewFrame(parser.NOP)
}

func processDefault(frame parser.Frame) parser.Frame {
	log.Debug("processDefault")
	return createErrorFrameWithMessage("Unknown Frame")
}
