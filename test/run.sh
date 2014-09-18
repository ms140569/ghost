#!/bin/bash

declare -A protocolTests

protocolTests=( 
    ["connect-and-send"]=0
	["connect-and-send-with-receipt"]=0
	["connect"]=0
	["connect-too-many-headers"]=1
	["connect-with-data"]=1
	["connect-without-both-mandatory"]=1
	["connect-header-key-too-long"]=1
	["connect-header-val-too-long"]=1
	["connect-disconnect"]=0
	["connect-disconnect-send"]=1
	["garbage10k"]=1   
	["garbage"]=1
	["stomp"]=0
	["send-with-receipt"]=1
	["stomp-with-data"]=1
	["header-corrupt1"]=1 
	["invalid-message"]=1
	["plain-connect"]=0
)

runTest() {
    # cat $1 |nc -q 1 localhost $PORT|grep -q ^$2 && return 0
    echo $1 $2
    ../ghostd -f $1 2>/dev/null 1>/dev/null
    return 0
}


for individualTest in "${!protocolTests[@]}" 
do 
	echo -n "$individualTest : "
	if (runTest $individualTest.stomp ${protocolTests["$individualTest"]})
	then
        
		echo OK
	else
		echo FAIL
	fi
	
done







