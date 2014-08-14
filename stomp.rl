/* -*-go-*-

Lexer for STOMP protocol, see http://stomp.github.io/

*/

package main

import (
	"log"
)

%%{
        machine lexer;
        write data;
}%%

func stomp_lexer(data string, tokenArray *[]Token) {
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
			    emitToken(Token{name: COMMAND, value: CommandForString(command)}, tokenArray) 
				};

			server_commands => {
			    command := data[ts:te];
				emitToken(Token{name: COMMAND, value: CommandForString(command)}, tokenArray) 
				};
		
		NULL =>   { emitToken(Token{name: NULL, value: nil}, tokenArray)};
		EOL =>    { emitToken(Token{name: EOL, value: nil}, tokenArray) };
		COLON =>  { emitToken(Token{name: COLON, value: nil}, tokenArray) };
		OCTET =>  { emitToken(Token{name: OCTET, value: nil}, tokenArray) };
		HEADER => { emitToken(Token{name: HEADER, value: data[ts:te]}, tokenArray) };
		STRING => { emitToken(Token{name: STRING, value: data[ts:te]}, tokenArray) };

		*|;

		write init;
		write exec;

	}%%

}

func Scanner(content string) []Token {
	log.Println("Scanner--------------------------------------");
    log.Println(string(content));

	tokenArray := []Token {}

	stomp_lexer(content, &tokenArray);

	log.Printf("Token appended, len: %d", len(tokenArray))

	return tokenArray;
}

func emitToken(token Token, tokenArray *[]Token) {
	// log.Printf("emitToken: %s", token)
	*tokenArray = append(*tokenArray, token)
}
