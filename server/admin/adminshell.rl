/* -*-go-*-

Handler for simple commands comming from a telnet connection:

dest list
dest create <name>
dest delete <name>
dest stat <name>

help
status
quit

*/

package admin

import (
	"github.com/ms140569/ghost/log"
)

%%{
        machine lexer;
        write data;
}%%

func command_lexer(data []byte, tokenArray *[]Shellcommand) {
	// These variables need to be predefined to get the ragel scanner running, 
    // see section 6.3 of the ragel userguide.
	act, ts, te, cs, p, pe := 0, 0, 0, 0, 0, len(data);

	var _, _, _ = act, ts, te; // This is to disable go's variable-declared-but-not-used error.

	%%{
		EOL = "\r"? . "\n";
		SPACE = " ";
		STRING = /[a-zA-Z0-9_\+\-\.\/\,]/+;

		commands = "status" | "help" | "quit" | "dest";

		main := |*

		commands => {
		    // command := data[ts:te];
		    emitToken(Shellcommand{}, tokenArray) 
			};
		
		EOL =>    { 
			emitToken(Shellcommand{}, tokenArray) 
		};


		*|;

		write init;
		write exec;

	}%%

}

func CommandScanner(content [] byte) []Shellcommand {
	log.Debug("CommandScanner--------------------------------------");
    log.Debug(string(content));

	tokenArray := []Shellcommand {}

	command_lexer(content, &tokenArray);

	log.Debug("Token appended, len: %d", len(tokenArray))

	return tokenArray;
}

func emitToken(token Shellcommand, tokenArray *[]Shellcommand) {
	*tokenArray = append(*tokenArray, token)
}
