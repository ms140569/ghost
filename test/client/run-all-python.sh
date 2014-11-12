#!/bin/bash

declare -A protocolTests

RED='\e[0;31m'
GREEN='\e[1;32m'
OFF='\e[0m'

protocolTests=( 
    ["connect"]=0
	["client"]=0
	["receipt"]=0
	["defect"]=1
)

for individualTest in "${!protocolTests[@]}" 
do 
	echo -n $individualTest.py ": "
    
    ./$individualTest.py 1>/dev/null 2>/dev/null

	if [ $? -eq ${protocolTests["$individualTest"]} ]
	then
        
		echo -e "${GREEN}SUCCESS${OFF}"
	else
		echo -e "${RED}FAIL${OFF}"
	fi
	
done
