package main

import (
	"flag"
	"github.com/ms140569/ghost/server"
)

var port = flag.Int("port", 7777, "TCP port to listen on")

func main() {
	flag.Parse()
	server.Server(*port)
}

