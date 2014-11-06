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

		log.Debug("Reading %d bytes", len(buffer))

		if err != nil {
			log.Fatal("Error reading file: %s", err)
			os.Exit(1)
		}

		parser.ParseFrames(buffer)
		os.Exit(0)
	}

	InitFrameQueue()

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
		log.Debug("Network read returned that much bytes:%d", n)
		buffer = buffer[0:n]

		if len(buffer) > 0 {
			bytesConsumed, frames, err := parser.ParseFrames(buffer)

			if err != nil { // log and/or process error but pass on valid frames
				log.Error(err.Error())
			}

			for _, frame := range frames {
				QueueFrame(conn, frame)
			}

			if bytesConsumed < len(buffer) { // reload new data and parse again
				log.Info("Parser did not consume all bytes.")
				log.Info("Received: %d, consumed by parsing: %d", len(buffer), bytesConsumed)
			}
		}

	}

	log.Debug("Connection from %v closed.", conn.RemoteAddr())

}
