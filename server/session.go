package server

import (
	"github.com/ms140569/ghost/log"
	"github.com/twinj/uuid"
	"net"
	"time"
)

type Session struct {
	Connection   net.Conn
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

func NewSession(conn net.Conn) Session {
	session := Session{isConnected: true, receivingHeartbeats: 0, sendingHeartbeats: 0, numberOfFramesReceived: 1}
	session.created = time.Now()
	session.Connection = conn

	// generate a session id
	session.id = uuid.NewV4().String()

	return session
}

func removeSessionFromMaps(conn net.Conn) {
	delete(sessions, conn)
	delete(sessionsToCheck, conn)
	delete(sessionsToKeepAlive, conn)
}

func (s *Session) Dump() {

	const layout = "Jan 2, 2006 at 3:04pm (MST)"

	log.Debug("Connection            : %o", s.Connection)
	log.Debug("isConnected           : %t", s.isConnected)
	log.Debug("isProtocol12          : %t", s.isProtocol12)
	log.Debug("id                    : %s", s.id)

	log.Debug("receivingHeartbeats   : %d", s.receivingHeartbeats)
	log.Debug("sendingHeartbeats     : %d", s.sendingHeartbeats)

	log.Debug("created               : %s seconds: %d", s.created.Format(layout), s.created.Second())
	log.Debug("lastKeepaliveReceived : %s seconds: %d", s.lastKeepaliveReceived.Format(layout), s.lastKeepaliveReceived.Second())

	log.Debug("numberOfFramesReceived: %d", s.numberOfFramesReceived)

}

func SessionStatus() {
	log.Debug("Number of sessions                  : %d", len(sessions))
	log.Debug("Number of sessions to check         : %d", len(sessionsToCheck))
	log.Debug("Number of sessions to send heartbeat: %d", len(sessionsToKeepAlive))
}
