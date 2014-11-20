package globals

import (
	"fmt"
	"github.com/ms140569/ghost/constants"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/storage"
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
	Port         int
	Testfilename string
	Testmode     bool
	Loglevel     string
	Storage      storage.Storekeeper
}

func NewConfig(flagBundle FlagBundle) {

	var configFile ConfigFile

	if len(flagBundle.Configfile) > 0 {
		configFile = readConfigFile(flagBundle.Configfile)
	}

	Config = Configuration{configFile.Port, configFile.Testfilename, configFile.Testmode, configFile.Loglevel, nil}

	// merging command line parameters

	if flagBundle.Port > 0 {
		Config.Port = flagBundle.Port
	}

	if len(flagBundle.Testfilename) > 0 {
		Config.Testfilename = flagBundle.Testfilename
	}

	Config.Testmode = flagBundle.Testmode

	if len(flagBundle.Loglevel) > 0 {
		Config.Loglevel = flagBundle.Loglevel
	}

	// checks and defaults

	if len(Config.Testfilename) > 0 {
		Config.Testmode = true
	}

	if Config.Port == 0 {
		Config.Port = 7777
	}

	log.SetSystemLogLevelFromString(Config.Loglevel)

	if len(configFile.Storage) > 0 {
		Config.Storage = storage.New(configFile.Storage)
	} else {
		Config.Storage = storage.New("mem:")
	}
}

func (c *Configuration) GetServerGreeting() string {
	return fmt.Sprintf(constants.GhostServerName+" version "+constants.GhostVersionNumber+" running on port: %s", strconv.Itoa(Config.Port))
}

func (c *Configuration) GetServerVersionString() string {
	return constants.GhostServerName + "/" + constants.GhostVersionNumber
}

func (c *Configuration) Dump() {
	log.Debug("Port: %d", c.Port)
	log.Debug("Testfilename: %s", c.Testfilename)
	log.Debug("Testmode: %t", c.Testmode)
}
