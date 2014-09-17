package main

import (
	"flag"
	"github.com/ms140569/ghost/server"
	"github.com/ms140569/ghost/globals"
)

func main() {
	var port = flag.Int("p", 7777, "TCP port to listen on")
	var filename = flag.String("f", "", "Filename of file containing frames.")
	var testmode = flag.Bool("t", false, "Testmode")

	flag.Parse()

	globals.NewConfig(*port, *filename, *testmode)

	server.Server()
}
