package server

import (
	"errors"
	"fmt"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"net"
	"strconv"
	"strings"
	"time"
)

// This goroutine sends heartbeats over connections, where this is needed.
func HeartBeatSender() {
	log.Info("Heartbeat sender started")
	for {
		time.Sleep(5 * time.Second)
	}
}

// This goroutine checks whether clients are still alive.
func HeartBeatChecker() {

	// this computation might be changed in the long run.
	// by now we wake up two times more often than the commited
	// response frequency
	var interval int = globals.HeartbeatsMinimalInterval / 2

	log.Info("Heartbeat checker started, wakeup frequency is %d", interval)
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)

		for _, session := range sessionsToCheck {

			timeDiff := timeDifferenceInMillis(time.Now(), session.lastKeepaliveReceived)
			if timeDiff > session.receivingHeartbeats {
				log.Debug("Cutting of session %s", session.id)
				writeAnswer(session.Connection, createErrorFrameWithMessage("Disconnecting session, no heartbeats received since: "+strconv.Itoa(timeDiff)))
			}
		}
	}
}

func timeDifferenceInMillis(now time.Time, past time.Time) int {
	return int(now.Sub(past).Nanoseconds() / 1000000)
}

func initializeHeartbeatingForConnection(frame parser.Frame, session *Session) (string, error) {
	heartbeatConfig := frame.GetHeader("heart-beat")

	log.Debug("Client requested heartbeating, this style:" + heartbeatConfig)

	out, in, err := parseHeartbeat(heartbeatConfig)

	if err != nil {
		return "", errors.New("Problem parsing heartbeat values: " + err.Error())
	}

	if out == 0 {
		log.Debug("Client says it can not send heartbeats")
	} else {
		(*session).receivingHeartbeats = max(out, globals.HeartbeatsMinimalInterval)

		sessionsToCheck[frame.Connection] = session
	}

	if in == 0 {
		log.Debug("Client does not want to receive heartbeats")
	} else {
		(*session).sendingHeartbeats = max(in, globals.HeartbeatsSendingInterval)

		sessionsToKeepAlive[frame.Connection] = session
	}

	log.Debug("Heartbeat setup for this session, receiving: %d, sending: %d", (*session).receivingHeartbeats, (*session).sendingHeartbeats)
	return fmt.Sprintf("heart-beat:%d,%d", globals.HeartbeatsMinimalInterval, globals.HeartbeatsSendingInterval), nil
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

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

/*
Handles incoming frames or heartbeats.
*/

func updateKeepaliveRecords(conn net.Conn) {

	log.Debug("Updating status for connection: %o", conn)

	SessionStatus()

	_, present := sessionsToCheck[conn]

	if present {
		log.Debug("Session found to keep alive.")

		sessionsToCheck[conn].Dump()

		sessionsToCheck[conn].numberOfFramesReceived = sessionsToCheck[conn].numberOfFramesReceived + 1
		sessionsToCheck[conn].lastKeepaliveReceived = time.Now()
	} else {
		log.Debug("Connection NOT found to be kept alive.")
	}
}
