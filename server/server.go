package server

import (
	"fmt"
	"github.com/ms140569/ghost/parser"
	"github.com/ms140569/ghost/constants"
	"log"
	"net"
	"strconv"
)

var GhostPortAsString string

func produceServerGreeting() string {
	return fmt.Sprintf(constants.GhostServerName+" version "+constants.GhostVersionNumber+" running on port: %s", GhostPortAsString)
}

func Server(listenPort int) {

	GhostPortAsString = strconv.Itoa(listenPort)

	fmt.Println(produceServerGreeting() + "\n")

	listener, err := net.Listen("tcp", ":"+GhostPortAsString)

	if err != nil {
		log.Fatal("Unable to Listen on port "+GhostPortAsString, err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Problem ", err)
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	log.Println("Connection Handler invoked")

	buffer := make([]byte, constants.DefaultBufferSize)

	conn.Write([]byte(produceServerGreeting() + "\n"))

	for {
		n, err := conn.Read(buffer)
		if err != nil || n == 0 {
			conn.Close()
			break
		}
		log.Printf("Read returned that much bytes:%d", n)
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

	log.Printf("Connection from %v closed.", conn.RemoteAddr())

}