import threading
import time


class ClientInteractor(threading.Thread):
    i=777
    def __init__(self, val) -> None:
        super().__init__()
        self.i= val
        print ("Initializnig new Client from ", self.i)

    def run(self, prefix) -> int:
        while(True):
            print (prefix, " Hello ",self.i)
            self.i = self.i + 1
            time.sleep(1)
            if self.i%20==0:
                return self.i