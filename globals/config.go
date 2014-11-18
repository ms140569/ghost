package globals

import (
	"encoding/json"
	"fmt"
	"github.com/ms140569/ghost/log"
	"os"
	"strconv"
)

// These are all the flags provided on the programs command line
// parsed by golang's flag package
type FlagBundle struct {
	Port         int
	Testfilename string
	Testmode     bool
	Loglevel     string
	Configfile   string
}

// Global variable used by the rest of the system
var Config Configuration

type Configuration struct {
	Port              int
	Testfilename      string
	Testmode          bool
	Loglevel          string
	GhostPortAsString string
	ServerGreeting    string
}

// simple config reader, no merging no overlays etc.
func NewConfig(flagBundle FlagBundle) {

	if len(flagBundle.Configfile) > 0 {
		Config = readConfigFile(flagBundle.Configfile)
	} else {
		Config = Configuration{flagBundle.Port, flagBundle.Testfilename, flagBundle.Testmode, "Debug", "", ""}
	}

	if len(Config.Testfilename) > 0 {
		Config.Testmode = true
	}

	Config.GhostPortAsString = strconv.Itoa(Config.Port)
	Config.ServerGreeting = produceServerGreeting(Config.GhostPortAsString)
	log.SetSystemLogLevelFromString(Config.Loglevel)
}

func produceServerGreeting(GhostPortAsString string) string {
	return fmt.Sprintf(GhostServerName+" version "+GhostVersionNumber+" running on port: %s", GhostPortAsString)
}

func GetServerVersionString() string {
	return GhostServerName + "/" + GhostVersionNumber
}

func readConfigFile(filename string) Configuration {
	log.Debug("Reading configuration from file: %s", filename)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal("Error reading file: %s", err.Error())
		os.Exit(4)
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)

	if err != nil {
		log.Fatal("error:", err)
		os.Exit(5)
	}

	configuration.Dump()

	return configuration

}

func (c *Configuration) Dump() {
	log.Debug("Port: %d", c.Port)
	log.Debug("Testfilename: %s", c.Testfilename)
	log.Debug("Testmode: %t", c.Testmode)
}
