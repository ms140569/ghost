#!/bin/bash
PORT=7777

declare -A protocolTests

protocolTests=( 
								["connect-and-send"]="CONNECTED"
								["connect-and-send-with-receipt"]="CONNECTED"
								["connect"]="CONNECTED"
								["connect-too-many-headers"]="ERROR"
								["connect-with-data"]="ERROR"
								["connect-without-both-mandatory"]="ERROR"
								["connect-header-key-too-long"]="ERROR"
								["connect-header-val-too-long"]="ERROR"
								["connect-disconnect"]="RECEIPT"
								["connect-disconnect-send"]="ERROR"
								["garbage10k"]="ERROR"   
								["garbage"]="ERROR"
								["stomp"]="CONNECTED"
								["send-with-receipt"]="ERROR"
								["stomp-with-data"]="ERROR"
								["header-corrupt1"]="ERROR" 
								["invalid-message"]="ERROR"
								["plain-connect"]="CONNECTED"
								)

connectAndGrep() {
        cat $1 |nc -q 1 localhost $PORT|grep -q ^$2 && return 0
        return 1
}

for individualTest in "${!protocolTests[@]}" 
do 
	echo -n "$individualTest : "
	if (connectAndGrep $individualTest.stomp ${protocolTests["$individualTest"]})
	then
		echo OK
	else
		echo FAIL
	fi
	
done







