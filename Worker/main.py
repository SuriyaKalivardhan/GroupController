from clientInteractor import ClientInteractor
from controllerInteractor import ControllerInteractor
import time
import logging

logging.basicConfig(format='%(asctime)s %(message)s',level=logging.INFO)

def main():
    logging.info("Init")
    ctrInteractor = ControllerInteractor()
    ctrInteractor.start()
    while ctrInteractor.client is None:
        time.sleep(1)
    logging.info("Client controller initialized")
    # cliInteractor = ctrInteractor.client
    #logging.INFO("Main got result from Client controller %s", cliInteractor.debug())

if __name__ == "__main__":
    main()