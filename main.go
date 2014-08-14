package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

var port = flag.Int("port", 7777, "TCP port to listen on")

const DefaultBufferSize int = 4096

const GhostVersionLine string = "Ghost v0.0"

var GhostPortAsString string

func produceServerGreeting() string {
	return fmt.Sprintf(GhostVersionLine+" running on port: %s", GhostPortAsString)
}

func main() {
	flag.Parse()
	GhostPortAsString = strconv.Itoa(*port)

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

	buffer := make([]byte, DefaultBufferSize)

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
			log.Printf("Getting that much data:%d", len(buffer))
			Scanner(string(buffer))
		}

	}

	log.Printf("Connection from %v closed.", conn.RemoteAddr())

}
