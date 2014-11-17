package storage

import (
/*	"github.com/ms140569/ghost/log"
	"github.com/twinj/uuid"
	"net"
	"time" */
)

// this is the most simple implementation of the storage
// interface: storing everythin in memory and loose it if
// one shuts down your process

type Memory struct{}

func (m Memory) Initialize() bool {
	return true
}
