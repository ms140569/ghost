#!/usr/bin/python
import logging, sys
logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)

from stompest.config import StompConfig
from stompest.sync import Stomp

CONFIG = StompConfig('tcp://localhost:7777', version='1.2')

if __name__ == '__main__':
    client = Stomp(CONFIG)
    client.connect()
    client.disconnect()

sys.exit(0)
