package globals

import (
	"fmt"
	"github.com/ms140569/ghost/log"
	"strconv"
)

// These are all the flags provided on the programs command line
// parsed by golang's flag package
type FlagBundle struct {
	Port       int
	Filename   string
	Testmode   bool
	Loglevel   string
	Configfile string
}

// Global variable used by the rest of the system
var Config Configuration

type Configuration struct {
	Port              int
	Filename          string
	Testmode          bool
	GhostPortAsString string
	ServerGreeting    string
}

func NewConfig(flagBundle FlagBundle) {

	Config = Configuration{flagBundle.Port, flagBundle.Filename, flagBundle.Testmode, "", ""}

	if len(flagBundle.Filename) > 0 {
		Config.Testmode = true
	}

	Config.GhostPortAsString = strconv.Itoa(flagBundle.Port)
	Config.ServerGreeting = produceServerGreeting(Config.GhostPortAsString)
	log.SetSystemLogLevelFromString(flagBundle.Loglevel)
}

func produceServerGreeting(GhostPortAsString string) string {
	return fmt.Sprintf(GhostServerName+" version "+GhostVersionNumber+" running on port: %s", GhostPortAsString)
}

func GetServerVersionString() string {
	return GhostServerName + "/" + GhostVersionNumber
}
