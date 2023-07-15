from .controller import Controller

from threading import Thread

class KeyboardReceiver(Thread):
    HELP = "P: pick up / H: hang up / Q: quit / 0-9: enter number / ?: Help"

    def __init__(self, controller: Controller):
        super().__init__()
        self.__controller: Controller = controller
        self.__do_quit: bool = False
    
    def run(self):
        self.__controller.start()
        print("Use your Keyboard to run the scenario")
        print(self.HELP)
        while not self.__do_quit:
            self.__prompt()

    def __prompt(self):
        selection = input("> ").lower()
        if selection == "p":
            print("Pick up phone")
            self.__controller.pick_up()
        elif selection == "h":
            print("Hang up phone")
            self.__controller.hang_up()
        elif selection == "?":
            print(KeyboardReceiver.HELP)
        elif selection == "q":
            print("Quit")
            self.__controller.quit()
            self.__do_quit = True
        elif str(selection).isdigit() and int(selection) > -1 and int(selection) <= 9:
            print(f"Dial number {selection}")
            self.__controller.dial(int(selection))
        else:
            print("Invalid command, type ? for help")


class MqttReceiver:
    def __init__(self, controller: Controller):
        self.__controller: Controller = controller
        self.__do_quit: bool = False

    def run(self):
        raise NotImplementedError()