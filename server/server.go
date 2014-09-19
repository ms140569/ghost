package server

import (
	"github.com/ms140569/ghost/globals"
	"github.com/ms140569/ghost/log"
	"github.com/ms140569/ghost/parser"
	"io/ioutil"
	"net"
	"os"
)

func Server() {

	if globals.Config.Testmode && len(globals.Config.Filename) > 0 {
		log.Info("Running in testmode, using file: %s", globals.Config.Filename)

		buffer, err := ioutil.ReadFile(globals.Config.Filename)

		if err != nil {
			log.Fatal("Error reading file: %s", err)
			os.Exit(1)
		}

		parser.ParseFrames(buffer)
		os.Exit(0)
	}

	log.Info(globals.Config.ServerGreeting + "\n")

	listener, err := net.Listen("tcp", ":"+globals.Config.GhostPortAsString)

	if err != nil {
		log.Fatal("Unable to Listen on port: %s"+globals.Config.GhostPortAsString, err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Error("Problem: %s", err)
			continue
		}

		go handleConnection(globals.Config.ServerGreeting, conn)
	}

}

func handleConnection(greeting string, conn net.Conn) {
	log.Debug("Connection Handler invoked")

	buffer := make([]byte, globals.DefaultBufferSize)

	conn.Write([]byte(greeting + "\n"))

	for {
		n, err := conn.Read(buffer)
		if err != nil || n == 0 {
			conn.Close()
			break
		}
		log.Debug("Read returned that much bytes:%d", n)
		buffer = buffer[0:n]

		n, err = conn.Write([]byte("Thanks for the data.\n"))

		if err != nil {
			conn.Close()
			break
		}

		if len(buffer) > 0 {
			parser.ParseFrames(buffer)
		}

	}

	log.Debug("Connection from %v closed.", conn.RemoteAddr())

}
