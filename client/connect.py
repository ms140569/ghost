#!/usr/bin/python
import logging
logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)

from stompest.config import StompConfig
from stompest.sync import Stomp

CONFIG = StompConfig('tcp://localhost:1405', version='1.2')

if __name__ == '__main__':
    client = Stomp(CONFIG)
    client.connect()
    client.disconnect()
