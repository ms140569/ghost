#!/bin/bash

declare -A protocolTests

RED='\e[0;31m'
GREEN='\e[1;32m'
OFF='\e[0m'

protocolTests=( 
["send-with-alpha"]=1     
["send-without"]=0                    
["send-with-positive-too-high"]=1      
["send-with-zero-no-data-no-nullbyte"]=1
["send-with-empty"]=1     
["send-with-positive-no-nullbyte"]=1  
["send-with-positive-too-low"]=1       
["send-with-zero-toomuch"]=1
["send-with-negative"]=1  
["send-with-positive"]=0              
["send-with-zero-data-no-nullbyte"]=1  
["send-with-zero-valid"]=0
["send-with-positive-multiple-nulls"]=0
)

runTest() {
    ../../ghostd -f $1 2>/dev/null 1>/dev/null
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
