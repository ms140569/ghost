package storage

import (
	"errors"
	"fmt"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"strings"
)

// This storage implementation is supposed to store and
// load all its data in Go's GOB serialization format.

type FileStorage struct {
	destinations []string
}

func NewFileStorage(configString string) FileStorage {

	storage := FileStorage{}

	if len(configString) > 0 {
		// parse given destinations give like this:
		// destination1=name&destination2=nextname

		for _, dest := range strings.Split(configString, "&") {
			if strings.HasPrefix(dest, "dest") {
				// arr := strings.Split(dest, "=")
				//log.Debug("Destination: %s", arr[1])
			}

		}
	}

	storage.Dump()

	return storage
}

func (m FileStorage) Initialize() error {
	return nil
}

func (m FileStorage) SendFrame(dest string, frame parser.Frame) error {
	log.Debug("Storing Frame at destination: %s", dest)
	return nil
}

func (m FileStorage) Subscribe(destination string, id string) error {
	log.Debug("Subscribe to destination %s using id: %s", destination, id)

	found := false

	for _, dest := range m.destinations {
		if dest == destination {
			found = true
		}
	}

	if !found {
		msg := fmt.Sprintf("Destination not found: %s", destination)
		log.Error(msg)
		return errors.New(msg)
	}

	log.Debug("Saving subscription with id %s", id)
	return nil
}

func (m FileStorage) Dump() {
	log.Debug("Memory storage provider dump:")
	for _, dest := range m.destinations {
		log.Debug("Destination: %s", dest)
	}
}
