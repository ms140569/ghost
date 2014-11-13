package server

import (
	"errors"
	"fmt"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
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
	log.Info("Heartbeat checker started")
	for {
		time.Sleep(5 * time.Second)
	}

}

func initializeHeartbeatingForConnection(frame parser.Frame, logicalConnection *LogicalConnection) (string, error) {
	heartbeatConfig := frame.GetHeader("heart-beat")

	log.Debug("Client requested heartbeating, this style:" + heartbeatConfig)

	out, in, err := parseHeartbeat(heartbeatConfig)

	if err != nil {
		return "", errors.New("Problem parsing heartbeat values: " + err.Error())
	}

	if out == 0 {
		log.Debug("Client says it can not send heartbeats")
	} else {
		(*logicalConnection).receivingHeartbeats = max(out, globals.HeartbeatsMinimalInterval)

		connectionsToCheck[frame.Connection] = *logicalConnection
	}

	if in == 0 {
		log.Debug("Client does not want to receive heartbeats")
	} else {
		(*logicalConnection).sendingHeartbeats = max(in, globals.HeartbeatsSendingInterval)

		connectionsToKeepAlive[frame.Connection] = *logicalConnection
	}

	log.Debug("Heartbeat setup for this session, receiving: %d, sending: %d", (*logicalConnection).receivingHeartbeats, (*logicalConnection).sendingHeartbeats)
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
