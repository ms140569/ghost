package storage

import (
	"errors"
	"fmt"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"strings"
)

// here goes the standard storage interfaces to be
// implemented by the various provider

type Storekeeper interface {
	Initialize() error
	SendFrame(destination string, frame parser.Frame) error
	Subscribe(destination string, id string) error

	// Destinations

	ListDestinations() []string
	CreateDestinattion(destination string) error
	DeleteDestination(destination string) error
	StatusDestination(destination string) (string, error)
}

/*

Creates a new storageprovider based on the given configuration
which looks like this:

Storage configuration
<storage-provider>:<provider specific information>

mem:field1=value1&field2=value2&field3=value3
mem:destination1=name&destination2=nextname
file:filename=data.ghostdata
mock:package.Implementation
mongo:host=localhost&username=testuser&password=gonzo
rocksdb:host=localhost&username=testuser&password=gonzo&port=6745
cockroachdb:
sqlite3:

*/

func New(configString string) (Storekeeper, error) {

	var method string
	var config string

	if strings.Contains(configString, ":") {
		arr := strings.Split(configString, ":")
		method = arr[0]
		config = arr[1]
	} else {
		method = configString
		config = ""
	}

	// log.Debug("Creating storage provider for config: %s", configString)

	switch method {
	case "mem":
		return NewMemory(config), nil
	case "file":
		return NewFileStorage(config), nil
	default:
		msg := fmt.Sprintf("Storage provider unkonwn: %s", method)
		log.Fatal(msg)
		return nil, errors.New(msg)

	}

	return nil, errors.New("No default storage provider. Supply configuration!")
}
