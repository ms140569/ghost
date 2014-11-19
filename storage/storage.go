package storage

import (
	"github.com/ms140569/ghost/log"
	"os"
	"strings"
)

// here goes the standard storage interfaces to be
// implemented by the various provider

type Storekeeper interface {
	Initialize() bool
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

func CreateStorageprovider(configString string) Storekeeper {
	arr := strings.Split(configString, ":")
	method := arr[0]
	// config := arr[1]

	log.Debug("Creating storage provider for config: %s", configString)

	switch method {
	case "mem":
		return Memory{}
	default:
		log.Fatal("Storage provider unkonwn: %s", method)
		os.Exit(5)
	}

	return Memory{}
}
