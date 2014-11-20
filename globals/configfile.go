package globals

import (
	"code.google.com/p/gcfg"
	"github.com/ms140569/ghost/log"
	"os"
)

type ConfigFile struct {
	Port         int
	Testfilename string
	Testmode     bool
	Loglevel     string
	Storage      string
}

func readConfigFile(filename string) ConfigFile {

	// We have to give a struct-in-struct here to match the ini-style sections
	type Config struct {
		Basic ConfigFile
	}

	cfg := Config{}
	err := gcfg.ReadFileInto(&cfg, filename)

	if err != nil {
		log.Fatal("Error reading file: %s", err.Error())
		os.Exit(4)
	}

	return cfg.Basic
}
