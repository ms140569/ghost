/* -*-go-*-

Handler for simple commands comming from a telnet connection:

dest list
dest create <name>
dest delete <name>
dest stat <name>

help
status
quit
show <name>

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

	cs, p, pe := 0, 0, len(data);

	%%{

		action emitDest {
			log.Debug("emmitting stuff: %s", string(data[0:pe]))
			log.Debug("cs, p, pe - %d, %d, %d", cs, p, pe)
		}

		action emitSimple {
			log.Debug("Emmitting simple command: %s", string(data[0:pe]))
			log.Debug("cs, p, pe - %d, %d, %d", cs, p, pe)
		}

		eol = "\r"? . "\n";

		simple_cmd = ("status" | "help" | "quit" | "show") > emitSimple;

		base_cmd_dest = "dest";
		sub_cmd_dest = "list" | "create" | "delete" | "stat";
		dest_grp = base_cmd_dest space sub_cmd_dest > emitDest;

		cmd = ( simple_cmd | dest_grp );

		lineval = cmd (space any+)?;
		line = lineval eol;

		main := line;

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
