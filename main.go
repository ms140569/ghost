package main

import (
	"flag"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/server"
	"github.com/ms140569/ghost/server/admin"
)

func main() {
	globals.NewConfig(parseFlags())

	admin.StartTelnetAdminServer()

	server.Server()
}

func parseFlags() globals.FlagBundle {
	flagBundle := globals.FlagBundle{}

	flag.IntVar(&flagBundle.Port, "p", 0, "TCP port to listen on")
	flag.StringVar(&flagBundle.Testfilename, "f", "", "Filename of file containing frames.")
	flag.BoolVar(&flagBundle.Testmode, "t", false, "Testmode")
	flag.StringVar(&flagBundle.Loglevel, "l", "", "Loglevel the server is running with")
	flag.StringVar(&flagBundle.Configfile, "c", "", "Program config file")

	flag.Parse()

	return flagBundle
}
