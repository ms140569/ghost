package storage

import (
//	"github.com/ms140569/ghost/log"
//	"github.com/twinj/uuid"
//	"net"
//	"time"
)

// here goes the standard storage interfaces to be
// implemented by the various provider

type Storekeeper interface {
	Initialize() bool
}
