/* -*-go-*-

Lexer for STOMP protocol, see http://stomp.github.io/

*/

package parser

import (
	"github.com/ms140569/ghost/log"
)

%%{
        machine lexer;
        write data;
}%%

func stomp_lexer(data []byte, tokenArray *[]Token) {
	// These variables need to be predefined to get the ragel scanner running, 
    // see section 6.3 of the ragel userguide.
	act, ts, te, cs, p, pe, eof := 0, 0, 0, 0, 0, len(data), len(data);

	var _, _, _ = act, ts, te; // This is to disable go's variable-declared-but-not-used error.

	numberOfEOLs := 0

	%%{
		EOL = "\r"? . "\n";
		COLON = ":";
		STRING = /[a-zA-Z0-9\+\-\.]/+;
		HEADER = STRING . COLON . STRING? . EOL; 

		client_commands = "SEND" | "SUBSCRIBE" | "UNSUBSCRIBE" | "BEGIN" | "COMMIT" | "ABORT" | "ACK" | "NACK" | "DISCONNECT" | "CONNECT" | "STOMP";
		server_commands = "CONNECTED" | "MESSAGE" | "RECEIPT" | "ERROR";
		all_commands = client_commands | server_commands;

		main := |*

		all_commands => {
		    command := data[ts:te];
		    emitToken(Token{name: COMMAND, value: CommandForString(string(command)), nextPos: te}, tokenArray) 
			};
		
		EOL =>    { 
			numberOfEOLs += 1
			emitToken(Token{name: EOL, value: nil, nextPos: te}, tokenArray) 
			if numberOfEOLs == 2 { // the second EOL marks the end of the header section. Data ought not to be slurped by the scanner.
				return
			}
		};

		HEADER => { emitToken(Token{name: HEADER, value: string(data[ts:te]), nextPos: te}, tokenArray) };

		*|;

		write init;
		write exec;

	}%%

}

func Scanner(content [] byte) []Token {
	log.Println("Scanner--------------------------------------");
    log.Println(string(content));

	tokenArray := []Token {}

	stomp_lexer(content, &tokenArray);

	log.Printf("Token appended, len: %d", len(tokenArray))

	return tokenArray;
}

func emitToken(token Token, tokenArray *[]Token) {
	*tokenArray = append(*tokenArray, token)
}
