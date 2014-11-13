package server

import (
	"github.com/ms140569/ghost/log"
	"net"
	"time"
)

type Session struct {
	isConnected  bool
	isProtocol12 bool
	id           string

	// in milliseconds
	receivingHeartbeats int
	sendingHeartbeats   int

	created               time.Time
	lastKeepaliveReceived time.Time

	// stats
	numberOfFramesReceived uint64
}

type SessionMap map[net.Conn]*Session

var sessions SessionMap = make(SessionMap)
var sessionsToCheck SessionMap = make(SessionMap)
var sessionsToKeepAlive SessionMap = make(SessionMap)

func (s *Session) Dump() {

	const layout = "Jan 2, 2006 at 3:04pm (MST)"

	log.Debug("isConnected           : %t", s.isConnected)
	log.Debug("isProtocol12          : %t", s.isProtocol12)
	log.Debug("id                    : %s", s.id)

	log.Debug("receivingHeartbeats   : %d", s.receivingHeartbeats)
	log.Debug("sendingHeartbeats     : %d", s.sendingHeartbeats)

	log.Debug("created               : %s", s.created.Format(layout))
	log.Debug("lastKeepaliveReceived : %s", s.lastKeepaliveReceived.Format(layout))

	log.Debug("numberOfFramesReceived: %d", s.numberOfFramesReceived)

}

func SessionStatus() {
	log.Debug("Number of sessions                  : %d", len(sessions))
	log.Debug("Number of sessions to check         : %d", len(sessionsToCheck))
	log.Debug("Number of sessions to send heartbeat: %d", len(sessionsToKeepAlive))
}
