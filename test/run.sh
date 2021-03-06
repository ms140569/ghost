#!/bin/bash

declare -A protocolTests

RED='\e[0;31m'
GREEN='\e[1;32m'
OFF='\e[0m'

protocolTests=( 
    ["connect-and-send"]=0
	["connect-and-send-with-receipt"]=0
	["connect"]=0
	["connect-too-many-headers"]=1
	["connect-with-data"]=1 # Only SEND Frames might have a body, see spec.
	# ["connect-without-both-mandatory"]=1 # stomp 1.2 clients ought to set accep-version and host header.
	["connect-header-key-too-long"]=1
	["connect-header-val-too-long"]=1
	["connect-disconnect"]=0
	["connect-disconnect-send"]=0
	["garbage10k"]=1   
	["garbage"]=1
	["stomp"]=0
	# ["send-with-receipt"]=1
	["stomp-with-data"]=1 # Only SEND Frames might have a body, see spec.
	["header-corrupt1"]=1 
	["invalid-message"]=1
	["plain-connect"]=0
)

runTest() {
    ../ghostd -f $1 2>/dev/null 1>/dev/null
    return $?
}

for individualTest in "${!protocolTests[@]}" 
do 
	echo -n $individualTest.stomp ": "
    
    runTest $individualTest.stomp

	if [ $? -eq ${protocolTests["$individualTest"]} ]
	then
        
		echo -e "${GREEN}SUCCESS${OFF}"
	else
		echo -e "${RED}FAIL${OFF}"
	fi
	
done
