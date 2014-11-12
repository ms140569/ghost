#!/usr/bin/python
import logging
logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)

from stompest.config import StompConfig
from stompest.sync import Stomp

CONFIG = StompConfig('tcp://localhost:7777', version='1.2')
QUEUE = '/queue/test'

if __name__ == '__main__':
    client = Stomp(CONFIG)
    client.connect(connectedTimeout=4)
    client.send(QUEUE, 'test message 1')
    client.send(QUEUE, 'test message 2')
    client.disconnect()
