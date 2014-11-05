#!/bin/bash

RED='\e[0;31m'
GREEN='\e[1;32m'
OFF='\e[0m'

runTest() {
    ../../ghostd -f $1 2>/dev/null 1>/dev/null
    return $?
}

for individualTest in `ls *.stomp` 
do 
	echo -n $individualTest ": "
    
    runTest $individualTest

	if [ $? -eq 0 ]
	then
		echo -e "${GREEN}SUCCESS${OFF}"
	else
		echo -e "${RED}FAIL${OFF}"
	fi
	
done
