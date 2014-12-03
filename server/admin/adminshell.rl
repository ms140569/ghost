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
	"strings"
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

	valid := false;

	%%{

		action mark { mark = p }
		
		action saveBase { 
			// log.Debug("saveBase() -> cs, p, pe - %d, %d, %d", cs, p, pe)
			baseCommand = string(data[mark:p]); 
			mark = p
		}

		action saveSub { 
			// log.Debug("saveSub -> cs, p, pe - %d, %d, %d", cs, p, pe)
			subCommand = string(data[mark:p]);
			mark = p
		}

		action saveParam { 
			// log.Debug("saveParam -> cs, p, pe - %d, %d, %d", cs, p, pe)
			param = strings.TrimSpace(string(data[mark:p]));
			mark = p
		}

		action validLine { 
			// log.Debug("VALID LINE");
			valid = true;
		}


		eol = "\r"? . "\n";

		# single commands

		single_cmd = ("status" | "help" | "quit") >mark %saveBase;
		single_with_param = "show" >mark %saveBase;

		fetch_single_param = (single_with_param space+ ^space+ ) >mark %saveParam;

		single_grp = ( single_cmd | fetch_single_param );

		# destination related commands

		base_cmd_dest = "dest" >mark %saveBase;

		sub_cmd_dest_single = "list" >mark %saveSub;
		sub_cmd_dest_param = "create" >mark %saveSub | "delete" >mark %saveSub | "stat" >mark %saveSub;

		fetch_dest_param = base_cmd_dest space+ sub_cmd_dest_param space+ ^space+ >mark %saveParam;
		dest_grp = ( base_cmd_dest space sub_cmd_dest_single | fetch_dest_param );

		cmd = ( single_grp | dest_grp ) %validLine;

		line = cmd eol;

		main := line;

		write init;
		write exec;

	}%%

    log.Debug("\n") 
    log.Debug("Command    : %s", baseCommand) 
	log.Debug("Subcommand : %s", subCommand) 
	log.Debug("Param      : %s", param) 

	if valid {
		return Shellcommand{name : ShellCommandNameForString(baseCommand), sub: subCommand, param: param}
	} else {
		return Shellcommand{name : UNDEF} 
	}

}

func CommandScanner(content [] byte) Shellcommand {
	log.Debug("CommandScanner--------------------------------------");
    log.Debug("Input      : %s", string(content));

	return command_lexer(content);

}
