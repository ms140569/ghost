package storage

import (
	"github.com/ms140569/ghost/log"
	"strings"
)

// this is the most simple implementation of the storage
// interface: storing everythin in memory and loose it if
// one shuts down your process

type Memory struct {
	destinations []string
}

func NewMemory(configString string) Memory {

	memory := Memory{}

	if len(configString) > 0 {
		// parse given destinations give like this:
		// destination1=name&destination2=nextname

		for _, dest := range strings.Split(configString, "&") {
			if strings.HasPrefix(dest, "dest") {
				arr := strings.Split(dest, "=")
				//log.Debug("Destination: %s", arr[1])
				memory.destinations = append(memory.destinations, arr[1])
			}

		}
	}

	memory.Dump()

	return memory
}

func (m Memory) Initialize() bool {
	return true
}

/*
func (m Memory) SendFrame(dest string, frame parser.Frame) {
	log.Debug("Storing Frame at destination: %s", dest)
}
*/

func (m Memory) Dump() {
	log.Debug("Memory storage provider dump:")
	for _, dest := range m.destinations {
		log.Debug("Destination: %s", dest)
	}
}
