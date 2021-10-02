from clientInteractor import ClientInteractor
from controllerInteractor import ControllerInteractor
import time

def main():
    ctrInteractor = ControllerInteractor(123)
    ctrInteractor.start()
    while ctrInteractor.client is None:
        time.sleep(1)
    cliInteractor = ctrInteractor.client
    print(cliInteractor.run("Main Prefif"))

if __name__ == "__main__":
    main()