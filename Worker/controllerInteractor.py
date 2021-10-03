import threading
import time
import logging
import redis_utils
import random

from clientInteractor import ClientInteractor

class ControllerInteractor(threading.Thread):
    client=None
    redisClient=None
    host="localhost"
    port=6379
    passwd=""
    redisClient=None
    pubsub=None
    id=None
    listenerChannel="None"
    controllerBootChannel = "ControllerBootChannel.v1"

    def __init__(self) -> None:
        id = random.randint(0,100)
        super().__init__()
        self.redisClient = redis_utils.connect_redis(self.host, self.port, self.passwd)
        self.pubsub = self.redisClient.pubsub()
        self.listenerChannel = format(id)+".listen"
        self.pubsub.subscribe(self.listenerChannel)
        self.redisClient.publish(self.controllerBootChannel, "ullala")

    def run(self):
        while(True):
            msg = self.pubsub.get_message(timeout=None)
            logging.info(msg)