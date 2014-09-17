package globals

import (
	"fmt"
	"strconv"
)

// Global variable used by the rest of the system
var Config Configuration

type Configuration struct {
	Port int
	Filename string
	Testmode bool
	GhostPortAsString string
	ServerGreeting string
}

func NewConfig(port int, filename string, testmode bool) {

	Config = Configuration{port, filename, testmode, "", ""}

	if len(filename) > 0 {
		Config.Testmode = true
	}

	Config.GhostPortAsString = strconv.Itoa(port)
	Config.ServerGreeting = produceServerGreeting(Config.GhostPortAsString)
}

func produceServerGreeting(GhostPortAsString string) string {
	return fmt.Sprintf(GhostServerName+" version "+GhostVersionNumber+" running on port: %s", GhostPortAsString)
}

