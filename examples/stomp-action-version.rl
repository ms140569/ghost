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

	cs, p, pe, eof := 0, 0, len(data), len(data);

	var _ = eof;


	%%{
		action mark {
			fmt.Println("Mark action");
		}

		action write_command {
			fmt.Println("Write command action");
		}
		action mark_key {
			fmt.Println("mark_key action");
		}

		action write_header {
			fmt.Println("Write header command action");
		}
		action finish_headers {
			fmt.Println("finish_headers action");
		}

		action consume_null {
			fmt.Println("Consume null command action");
		}

		action consume_octet {
			fmt.Println("consume octet action");
		}

		action write_body {
			fmt.Println("Write body command action");
		}

		action mark_frame {
			fmt.Println("Mark_frame action");
		}

		action finish_frame {
			fmt.Println("finish_frame command action");
		}

		action check_frame_size {
			fmt.Println("check_frame size command action");
		}

		action say_hello {
			fmt.Println("HELLO");
		}


		NULL = "\0";
		EOL = "\r"? . "\n";
		OCTET = any;

		client_commands = "SEND" | "SUBSCRIBE" | "UNSUBSCRIBE" | "BEGIN" | "COMMIT" | "ABORT" | "ACK" | "NACK" | "DISCONNECT" | "CONNECT" | "STOMP";
		server_commands = "CONNECTED" | "MESSAGE" | "RECEIPT" | "ERROR";
		command = (client_commands | server_commands) > mark % write_command . EOL;

		HEADER_ESCAPE = "\\" . ("\\" | "n" | "r" | "c");
		HEADER_OCTET = HEADER_ESCAPE | (OCTET - "\r" - "\n" - "\\" - ":");
		header_key = HEADER_OCTET+ > mark % mark_key;
		header_value = HEADER_OCTET* > mark;
		header = header_key . ":" . header_value;
		headers = (header % write_header . EOL)* % finish_headers . EOL;

		# consume_body = (NULL when consume_null | ^NULL when consume_octet)*;
		# body = consume_body >from(mark) % write_body <: NULL;

		body = OCTET* > say_hello;

		frame = ((command > mark_frame) :> headers :> (body @ finish_frame)) $ check_frame_size;

		stream := (EOL | frame)*;

		# stream := (EOL | command)*;

		write init;
		write exec;

	}%%

}

func Scanner(content string) {

     fmt.Println("Scanner");

     fmt.Println(string(content));

	stomp_lexer(content);

}
