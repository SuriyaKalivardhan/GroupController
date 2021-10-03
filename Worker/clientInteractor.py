import os
import threading
import time
import logging

import redis_utils 
import json


class ClientInteractor(threading.Thread):
    redisClient=None
    pubsub=None
    id=None
    clientId=None
    listenerChannel=None
    clientRegisterChannel = None

    def __init__(self, id, host, port, passwd, registerChannel, clientId) -> None:
        super().__init__()
        self.id = id
        self.clientId = clientId
        self.redisClient = redis_utils.connect_redis(host, port, passwd)
        self.pubsub = self.redisClient.pubsub()
        self.listenerChannel = format(self.id)+".model.listen"
        self.pubsub.subscribe(self.listenerChannel)
        request = {
            "id": format(self.id),
            "listenerChannel": self.listenerChannel,
            "method":"Register"
        }
        self.clientRegisterChannel = registerChannel
        self.redisClient.publish(self.clientRegisterChannel, json.dumps(request))
        logging.info("Initializnig new Client from %d", self.id)
        threading.Thread(target=self.monitor_client, args=(), daemon=True).start()

    def run(self):
        while(True):
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
        os._exit(3)
    
    def monitor_client(self):
        while True:
            clientSubs = self.redisClient.pubsub_numsub(self.clientRegisterChannel)
            clientListening = False

            for subsribers in clientSubs:
                if subsribers[0].decode('UTF-8') == self.clientRegisterChannel:
                    if subsribers[1] > 0:
                        clientListening = True

            if clientListening == False:
                logging.info("Crashing since no active clients to their channel")
                os._exit(3)

            time.sleep(1)
