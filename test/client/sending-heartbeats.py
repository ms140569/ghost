#!/usr/bin/python

import logging, sys
import time
from stompest.config import StompConfig
from stompest.sync import Stomp

logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)

CONFIG = StompConfig('tcp://localhost:7777', version='1.2')


def toTime(input):
    return time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(input))

if __name__ == '__main__':

    client = Stomp(CONFIG)

    client.connect(connectedTimeout=4, heartBeats=(2000,2500))

    print "--------------------------------------------------------------------"
    print "state        : ", client.session.state
    print "client HB    : ", client.session.clientHeartBeat
    print "server HB    : ", client.session.serverHeartBeat
    print "server       : ", client.session.server
    print "id           : ", client.session.id
    print "lastSent     : ", toTime(client.session.lastSent)
    print "lastReceived : ", toTime(client.session.lastReceived)
    
    start = time.time()
    elapsed = lambda t = None: (t or time.time()) - start

    times = lambda: 'elapsed: %.2f, last received: %.2f, last sent: %.2f' % (
        elapsed(), elapsed(client.lastReceived), elapsed(client.lastSent)
        )

    for x in range(0, 7):
        client.beat()
        print times()
        time.sleep(1)


#    while elapsed() < 2 * client.clientHeartBeat / 1000.0:
#        client.canRead(0.8 * client.serverHeartBeat / 1000.0) # poll server heart-beats
#        client.beat() # send client heart-beats
#        print times()

#    client.canRead()
#    print times()

    client.disconnect()
    sys.exit(0)
    
