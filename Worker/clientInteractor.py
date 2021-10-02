import threading
import time


class ClientInteractor(threading.Thread):
    i=777
    def __init__(self, val) -> None:
        super().__init__()
        i= val
        print ("Initializnig new Client from ", i)

    def run(self):
        while(True):
            print ("Hello ",self.i)
            self.i = self.i + 1
            time.sleep(5)