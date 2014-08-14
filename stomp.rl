/* -*-go-*-

Lexer for STOMP protocol, see http://stomp.github.io/

   ragel -Z -T0 -o stomp.go stomp.rl
   go build -o stomp stomp.go
   ./stomp

*/

package main

import (
        // "os"
        "fmt"
	// "io/ioutil"
)

%%{
        machine lexer;
        write data;
}%%

func stomp_lexer(data string) {
	// These variables need to be predefined to get the ragel scanner running, 
    // see section 6.3 of the ragel userguide.
	act, ts, te, cs, p, pe, eof := 0, 0, 0, 0, 0, len(data), len(data);

	var _, _, _ = act, ts, te; // This is to disable go's variable-declared-but-not-used error.

	%%{
		NULL = "\0";
		EOL = "\r"? . "\n";
		
		COLON = ":";
		STRING = /[a-zA-Z0-9\+\-\.]/+;
		HEADER = STRING . COLON; 
		OCTET = any;

		client_commands = "SEND" | "SUBSCRIBE" | "UNSUBSCRIBE" | "BEGIN" | "COMMIT" | "ABORT" | "ACK" | "NACK" | "DISCONNECT" | "CONNECT" | "STOMP";
		server_commands = "CONNECTED" | "MESSAGE" | "RECEIPT" | "ERROR";

		action mark {
			fmt.Println("Mark action");
		}

		action write_command {
			fmt.Println("Write command action");
		}

		main := |*

			client_commands => {
			    command := data[ts:te];
			    fmt.Println("Found client command: " + command); 
			    fmt.Println("CMD:", CommandForString(command));
				};

		    server_commands => { fmt.Println("Found server command: ") };
		
			NULL => { fmt.Println("Found NULL") };
		    EOL => { fmt.Println("Found EOL") };
		    COLON => { fmt.Println("Found COLON") };
		    OCTET => { fmt.Println("Found OCTET") };
		    HEADER => { fmt.Println("Found Header: " + data[ts:te]) };
		    STRING => { fmt.Println("Found String: " + data[ts:te]) };

		*|;

		write init;
		write exec;

	}%%

}

func Scanner(content string) {

     fmt.Println("Scanner");

     fmt.Println(string(content));

	stomp_lexer(content);

}
