from clientInteractor import ClientInteractor
from controllerInteractor import ControllreInteractor


def main():
    ctrInteractor = ControllreInteractor(123)
    ctrInteractor.start()

if __name__ == "__main__":
    main()