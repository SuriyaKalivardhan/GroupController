import threading
import time
import logging

import redis_utils 
import json


class ClientInteractor(threading.Thread):
    redisClient=None
    pubsub=None
    id=None
    listenerChannel=None
    modelRegisterchannel = None
    done=False

    def __init__():
        super.__init__()

    def __init__(self, id, host, port, passwd, registerChannel) -> None:
        super().__init__()
        self.id = id
        self.redisClient = redis_utils.connect_redis(host, port, passwd)
        self.pubsub = self.redisClient.pubsub()
        self.listenerChannel = format(self.id)+".model.listen"
        self.pubsub.subscribe(self.listenerChannel)
        request = {
            "id": format(self.id),
            "listenerChannel": self.listenerChannel,
            "method":"Register"
        }
        self.modelRegisterchannel = registerChannel
        self.redisClient.publish(self.modelRegisterchannel, json.dumps(request))
        logging.info("Initializnig new Client from %d", self.id)

    def run(self):
        while(self.done == False):
            msg = self.pubsub.get_message(timeout=3)
            if msg == None or msg["type"] != "message":
                continue
            request = json.loads(msg["data"])
            logging.info("Recived the data message %s from Model controller", request)
    
    def debug(self) -> str:
        return "Currently listening to  " + self.listenerChannel

    def UnRegister(self):        
        request = {
            "id": format(self.id),
            "listenerChannel": self.listenerChannel,
            "method": "Shutdown"
        }
        self.redisClient.publish(self.modelRegisterchannel, json.dumps(request))
        logging.info("Shutting down Client...")
        self.done = True