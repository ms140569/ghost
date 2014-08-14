/* -*-go-*-

Lexer for STOMP protocol, see http://stomp.github.io/

*/

package main

import (
	"log"
	"fmt"
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

		main := |*

			client_commands => {
			    command := data[ts:te];
			        emitToken(Token{name: COMMAND, value: CommandForString(command)}) 
				};

			server_commands => {
			    command := data[ts:te];
			        emitToken(Token{name: COMMAND, value: CommandForString(command)}) 
				};

		
		    NULL => { emitToken(Token{name: NULL, value: nil})};
     		EOL => { emitToken(Token{name: EOL, value: nil}) };
	    	COLON => { emitToken(Token{name: COLON, value: nil}) };
		    OCTET => { emitToken(Token{name: OCTET, value: nil}) };
		    HEADER => { emitToken(Token{name: HEADER, value: data[ts:te]}) };
		    STRING => { emitToken(Token{name: STRING, value: data[ts:te]}) };

		*|;

		write init;
		write exec;

	}%%

}

func Scanner(content string) {
	fmt.Println("Scanner--------------------------------------");
    fmt.Println(string(content));

	stomp_lexer(content);

}

func emitToken(token Token) {
	log.Printf("emitToken: %s", token)
}
