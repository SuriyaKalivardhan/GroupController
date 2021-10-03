import threading
import time
import logging

from redis import client
import redis_utils
import json
import random
import os

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
        self.id = random.randint(0,100)
        super().__init__()
        self.redisClient = redis_utils.connect_redis(self.host, self.port, self.passwd)
        self.pubsub = self.redisClient.pubsub()
        self.listenerChannel = format(self.id)+".listen"
        self.pubsub.subscribe(self.listenerChannel)
        threading.Thread(target=self.relay_to_group, args=(), daemon=True).start()

    def run(self):
        while(True):
            msg = self.pubsub.get_message(timeout=None)
            if msg == None or msg["type"] != "message":
                continue
            request = json.loads(msg["data"])
            if request["method"] == "Register":
                logging.info("Recived the register message %s", request)
                self.client = ClientInteractor(self.id,
                    request["redisHost"], request["redisPort"], request["redisPassword"],
                    request["registerChannel"], request["id"])
                self.client.start()
            elif request["method"] == "UnRegister":
                if self.client is not None:
                    self.client.UnRegister()
                logging.info("Shutting down controllrer inteactor.")
                os._exit(3)
    
    def relay_to_group(self):
        while True:

            request = None
            if self.client is not None:
                request = {
                "id": format(self.id),
                "listenerChannel": self.listenerChannel,
                "client": self.client.clientId
                }
            else:
                request = {
                    "id": format(self.id),
                    "listenerChannel": self.listenerChannel
                }
            self.redisClient.publish(self.controllerBootChannel, json.dumps(request))

            time.sleep(1)
