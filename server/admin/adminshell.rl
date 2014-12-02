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

func command_lexer(data []byte) Shellcommand {

	cs, p, pe := 0, 0, len(data);
	mark := 0;

	baseCommand := "";
	subCommand := "";
	param := "";

	%%{

		action mark { mark = p }
		
		action saveBase { 
			// log.Debug("SAVING: cs, p, pe - %d, %d, %d", cs, p, pe)
			baseCommand = string(data[mark:p]); 
			mark = p
		}

		action saveSub { 
			// log.Debug("SAVING: cs, p, pe - %d, %d, %d", cs, p, pe)
			subCommand = string(data[mark:p]); 
		}

		action saveParam { 
			// log.Debug("SAVING: cs, p, pe - %d, %d, %d", cs, p, pe)
			param = string(data[mark:p]); 
		}


		action emitDest {
			log.Debug("emmitting stuff: %s", string(data[0:pe]))
			log.Debug("cs, p, pe - %d, %d, %d", cs, p, pe)
		}

		action emitSimple {
			log.Debug("Emmitting simple command: %s", string(data[0:pe]))
			log.Debug("cs, p, pe - %d, %d, %d", cs, p, pe)
		}

		eol = "\r"? . "\n";

		simple_cmd = ("status" | "help" | "quit" | "show") >mark %saveBase;

		base_cmd_dest = "dest" >mark %saveBase;
		sub_cmd_dest = "list" | "create" | "delete" | "stat" %saveSub;

		dest_grp = base_cmd_dest space sub_cmd_dest;

		cmd = ( simple_cmd | dest_grp ) >mark;

		lineval = cmd (space any+)? %saveParam;
		line = lineval eol;

		main := line;

		write init;
		write exec;

	}%%

    log.Debug("Command    : %s", baseCommand) 
	log.Debug("Subcommand : %s", subCommand) 
	log.Debug("Param      : %s", param) 

	return Shellcommand{name : QUIT, sub: subCommand, param: param}

}

func CommandScanner(content [] byte) Shellcommand {
	log.Debug("CommandScanner--------------------------------------");
    log.Debug("Input      : %s", string(content));

	return command_lexer(content);

}
