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

		bytesConsumed, frames, err := parser.ParseFrames(buffer)

		log.Debug("Bytes consumed by parser: %d", bytesConsumed)

		if err != nil { // log and/or process error but pass on valid frames
			log.Error(err.Error())
			os.Exit(1)
		}

		// in Testmode we don't run async by queue-in the frames but process them
		// synchronously.

		for _, frame := range frames {
			ProcessFrame(frame)
		}

		os.Exit(0)
	}

	InitFrameQueues()

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

	for {
		buffer := make([]byte, globals.DefaultBufferSize)

		n, err := conn.Read(buffer)

		if err != nil || n == 0 {
			conn.Close()
			break
		}

		log.Debug("Network read returned that much bytes:%d", n)

		// if we've received exacty one byte which is an EOL
		// we've got a heartbeat

		if n == 1 && buffer[0] == 0x0A {
			log.Debug("It's a heartbeat coming from this connection: %o", conn)

			// FIXME: here we have to deal with incoming heartbeats
			continue
		}

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
