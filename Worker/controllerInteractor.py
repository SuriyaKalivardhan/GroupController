import threading
import time

from clientInteractor import ClientInteractor

class ControllreInteractor(threading.Thread):
    i=10
    client=None
    def __init__(self, inp) -> None:
        super().__init__()
        self.i=inp

    def run(self):
        while(True):
            print ("Current i ", self.i)
            time.sleep(1)
            self.i = self.i + 1
            if self.client is None and (self.i % 10== 0):
                self.client = ClientInteractor(self.i)
                self.client.start()