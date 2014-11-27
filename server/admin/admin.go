package admin

import (
	"github.com/ms140569/ghost/constants"
	"github.com/ms140569/ghost/log"
	"net"
	"os"
	"strconv"
)

func ListDestinations() []string {
	return make([]string, 3, 3)
}

func CreateDestinattion(destination string) {}

func DeleteDestination(destination string) {}

func StatusDestination(destination string) string {
	return "nix"
}

func StartTelnetAdminServer() {

	portAsString := strconv.Itoa(constants.DefaultTelnetPortNumber)

	listener, err := net.Listen("tcp", ":"+portAsString)

	if err != nil {
		log.Fatal("Unable to Listen on port: %s"+portAsString, err)
		os.Exit(1)
	}

	log.Info("Admin server running on port: %s", portAsString)

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				log.Error("Problem: %s", err)
				continue
			}

			go handleTelnetConnection(conn)
		}
	}()

}

func handleTelnetConnection(conn net.Conn) {
	log.Debug("Connection Handler invoked")

	for {
		buffer := make([]byte, constants.DefaultBufferSize)

		n, err := conn.Read(buffer)

		if err != nil || n == 0 {
			conn.Close()
			break
		}

		log.Debug("Network read returned that much bytes:%d", n)

		buffer = buffer[0:n]

		if len(buffer) > 0 {
			log.Debug("Go willy, go")

			_, err := conn.Write([]byte("What do you want?"))

			if err != nil {
				conn.Close()
			}

			conn.Close()
		}

	}

	log.Debug("Connection from %v closed.", conn.RemoteAddr())

}
