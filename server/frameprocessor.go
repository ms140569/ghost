package server

import (
	"errors"
	"fmt"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"github.com/twinj/uuid"
	"net"
	"os"
	"strconv"
	"strings"
)

type LogicalConnection struct {
	isConnected         bool
	isProtocol12        bool
	receivingHeartbeats int // in millis
	sendingHeartbeats   int // in millis
}

var connections map[net.Conn]LogicalConnection = make(map[net.Conn]LogicalConnection)

var inboundFrameQueue chan parser.Frame
var outboundFrameQueue chan parser.Frame

func InitFrameQueues() {

	inboundFrameQueue = make(chan parser.Frame, globals.QueueSizeInbound)
	outboundFrameQueue = make(chan parser.Frame, globals.QueueSizeOutbound)

	go FetchFrame()
	go FrameWriter()
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

		// we either send a full-fledged frame or a heartbeat

		var output []byte

		if frame.Command == parser.HEARTBEAT {
			output = []byte("\x0A")
		} else {
			output = []byte(frame.Render())
		}

		_, err := frame.Connection.Write(output)

		if err != nil {
			frame.Connection.Close()
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

	_, present := connections[frame.Connection]

	var answer parser.Frame

	if present {
		log.Debug("Connection know to server")
		answer = createErrorFrameWithMessage("already connected")
	} else {

		// FIXME: here should go the code to check for login and passcode

		if frame.HasHeader("login") && frame.HasHeader("passcode") {
			msg := "login and passcode not supported yet."
			return createErrorFrameWithMessage(msg)
		}

		// store connection for further reuse

		log.Debug("New connection, adding to map.")

		logicalConnection := LogicalConnection{isConnected: true, receivingHeartbeats: 0, sendingHeartbeats: 0}

		connections[frame.Connection] = logicalConnection

		answer = parser.NewFrame(parser.CONNECTED)
		answer.AddHeader("version:1.2")

		// Analyze and setup heartbeating ...

		if frame.HasHeader("heart-beat") {

			heartbeatConfig := frame.GetHeader("heart-beat")

			log.Debug("Client requested heartbeating, this style:" + heartbeatConfig)

			out, in, err := parseHeartbeat(heartbeatConfig)

			if err != nil {
				return createErrorFrameWithMessage("Problem parsing heartbeat values.")
			}

			if out == 0 {
				log.Debug("Client says it can not send heartbeats")
			} else {
				logicalConnection.receivingHeartbeats = max(out, globals.HeartbeatsMinimalInterval)
			}

			if in == 0 {
				log.Debug("Client does not want to receive heartbeats")
			} else {
				logicalConnection.sendingHeartbeats = max(in, globals.HeartbeatsSendingInterval)
			}

			heartbeatAnswer := fmt.Sprintf("heart-beat:%d,%d", globals.HeartbeatsMinimalInterval, globals.HeartbeatsSendingInterval)

			answer.AddHeader(heartbeatAnswer)

			log.Debug("Heartbeat session for this session, receiving: %d, sending: %d", logicalConnection.receivingHeartbeats, logicalConnection.sendingHeartbeats)

		}

		// generate a session id
		sessionId := uuid.NewV4()
		answer.AddHeader("session:" + sessionId.String())

		// server versioin
		answer.AddHeader("server:" + globals.GetServerVersionString())
	}

	return answer
}

/*
   Parsing the heartbeat header which is in the form of:
   outgoing,incoming
*/

func parseHeartbeat(s string) (int, int, error) {
	if strings.Count(s, ",") != 1 {
		return -1, -1, errors.New("Wrong number of commas in heartbeat header.")

	}

	arr := strings.Split(s, ",")

	outString := arr[0]
	inString := arr[1]

	if len(outString) == 0 || len(inString) == 0 {
		return -1, -1, errors.New("Either incoming or outgoing time not supplied.")
	}

	outVal, err := strconv.Atoi(outString)

	if err != nil {
		return -1, -1, errors.New("Error parsing outvalue")
	}

	inVal, err := strconv.Atoi(inString)

	if err != nil {
		return -1, -1, errors.New("Error parsing invalue")
	}

	if outVal < 0 || inVal < 0 {
		return -1, -1, errors.New("No negative values allowed")
	}

	return outVal, inVal, nil

}

func processSend(frame parser.Frame) parser.Frame {
	log.Debug("processSend")
	return parser.NewFrame(parser.NOP)
}

func processDisconnect(frame parser.Frame) parser.Frame {
	log.Debug("processDisconnect")
	return parser.NewFrame(parser.NOP)
}

func processDefault(frame parser.Frame) parser.Frame {
	log.Debug("processDefault")
	return createErrorFrameWithMessage("Unknown Frame")
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
