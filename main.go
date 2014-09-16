package main

import (
	"flag"
	"github.com/ms140569/ghost/server"
)

func main() {
	var port = flag.Int("p", 7777, "TCP port to listen on")
	var filename = flag.String("t", "", "Filename for tests")

	flag.Parse()

	config := server.NewConfig(*port, *filename)

	server.Server(config)
}
