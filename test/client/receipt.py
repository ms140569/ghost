#!/usr/bin/python

import logging, sys
logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)

from stompest.config import StompConfig
from stompest.sync import Stomp

CONFIG = StompConfig('tcp://localhost:7777', version='1.2')

if __name__ == '__main__':
    client = Stomp(CONFIG)
    client.connect(connectedTimeout=4)

    ID = 'bingo'

    client.send('/my/test/destination',       # destination 
                'THIS-IS-A-SEND-BODY',        # body
                headers={'receipt': ID})      # headers

    answer = client.receiveFrame()

    receiptID = answer.headers["receipt-id"] 
    returnValue = 0

    if receiptID != ID:
        print "Receipt header wrong:" + receiptID 
        returnValue = 1
    else:
        print "Correct receipt id received: " + receiptID

    client.disconnect()
    sys.exit(returnValue)
    

