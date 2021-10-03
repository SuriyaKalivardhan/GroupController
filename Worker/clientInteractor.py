import threading
import time
import logging


class ClientInteractor(threading.Thread):
    i=777
    def __init__(self, val) -> None:
        super().__init__()
        self.i= val
        logging.info("Initializnig new Client from %d", self.i)

    def run(self, prefix) -> int:
        while(True):
            logging.info("%s Hello %d", prefix, self.i)
            self.i = self.i + 1
            time.sleep(1)
            if self.i%20==0:
                return self.i