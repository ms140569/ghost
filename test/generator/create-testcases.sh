#!/bin/bash

RANGE='0 2'

printCommand() {
    echo SEND
}

printHeader() {

    echo destination:nirvana

    for i in `seq 1 $1` 
    do
        echo header$i:value$i
    done
}

printBody() {
    for i in `seq 1 $1` 
    do
        echo -n B
    done
}

printTrailingEOF() {
    for i in `seq 1 $1` 
    do
        echo
    done
}

printFrame() {
    printCommand
    printHeader $1
    echo 
    printBody $2
    printf '\00'
    printTrailingEOF $3
}

PREFIX=tc-conn

for a in `seq $RANGE` 
do
	for b in `seq $RANGE` 
    do
	    for c in `seq $RANGE` 
        do
            # echo $PREFIX-$a-$b-$c.stomp
            printFrame $a $b $c > $PREFIX-$a-$b-$c.stomp
            # echo ---------------------------------
	        
        done

    done

done

