package main

import (
	"flag"
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/server"
)

func main() {
	var port = flag.Int("p", 7777, "TCP port to listen on")
	var filename = flag.String("f", "", "Filename of file containing frames.")
	var testmode = flag.Bool("t", false, "Testmode")
	var logLevel = flag.String("l", "Debug", "Loglevel the server is running with")

	flag.Parse()

	globals.NewConfig(*port, *filename, *testmode, *logLevel)

	server.Server()
}
