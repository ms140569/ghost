package admin

import (
	"github.com/ms140569/ghost/constants"
	"github.com/ms140569/ghost/log"
	"net"
	"os"
	"strconv"
)

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

/*
Handler for simple commands like:

dest list
dest create <name>
dest delete <name>
dest stat <name>

help
status
quit


*/

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
			cmd := CommandScanner(buffer)
			runner := CommandRunner{conn: conn, cmd: cmd}
			runner.Execute()
		}

	}

	log.Debug("Connection from %v closed.", conn.RemoteAddr())

}

type CommandRunner struct {
	conn net.Conn
	cmd  Shellcommand
}

func (cr CommandRunner) Execute() {

	switch cr.cmd.name {
	case HELP:
		cr.reply(HelpForAllCommands())
	case STATUS:
		cr.reply("STATUS")
	case QUIT:
		cr.reply("QUIT")
		cr.conn.Close()
	case SHOW:
		cr.reply("SHOW")
	case DEST:
		switch cr.cmd.sub {
		case "list":
			cr.ListDestinations()
		case "create":
			cr.CreateDestination()
		case "delete":
			cr.DeleteDestination()
		case "stat":
			cr.StatusDestination()
		}

	case UNDEF:
		cr.reply("I do not understand.")
	default:
		cr.reply("I do not understand.")

	}

}

func (cr CommandRunner) reply(msg string) {
	_, err := cr.conn.Write([]byte(msg + "\n"))

	if err != nil {
		cr.conn.Close()
	}

}

func (cr CommandRunner) ListDestinations() {
	cr.reply("Listing destinations:")

	cr.reply(StringArrayPrinter(ListDestinations()))

}

func (cr CommandRunner) CreateDestination() {
	err := CreateDestination(cr.cmd.param)

	if err != nil {
		cr.reply(err.Error())
	} else {
		cr.reply("Done.")
	}

}

func (cr CommandRunner) DeleteDestination() {
	err := DeleteDestination(cr.cmd.param)

	if err != nil {
		cr.reply(err.Error())
	} else {
		cr.reply("Done.")
	}

}

func (cr CommandRunner) StatusDestination() {
	str, err := StatusDestination(cr.cmd.param)

	if err != nil {
		cr.reply(err.Error())
	} else {
		cr.reply(str)
	}

}

func StringArrayPrinter(arr []string) string {
	retVal := ""

	for _, val := range arr {
		retVal = retVal + val + "\n"
	}

	return retVal

}
