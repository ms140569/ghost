package server

import (
	"fmt"
	"strconv"
	"github.com/ms140569/ghost/constants"
)

type Config struct {
	Port int
	Filename string
	Testmode bool
	GhostPortAsString string
	ServerGreeting string
}

func NewConfig(port int, filename string) Config {

	var config Config = Config{port, filename, false, "", ""}

	if len(filename) > 0 {
		config.Testmode = true
	}

	config.GhostPortAsString = strconv.Itoa(port)
	config.ServerGreeting = produceServerGreeting(config.GhostPortAsString)

	return config
}

func produceServerGreeting(GhostPortAsString string) string {
	return fmt.Sprintf(constants.GhostServerName+" version "+constants.GhostVersionNumber+" running on port: %s", GhostPortAsString)
}

