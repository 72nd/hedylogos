from typing import Optional
from .controller import Controller

from threading import Thread

from readchar import readkey, key


class KeyboardReceiver(Thread):
    """
    Provides a simple shell for entering the different commands to the system. This receiver
    is mainly used for debugging. This mode can also be used with an numeric-keyboard to provide
    a simple user facing control panel without the need of setting up a real old dial phone.
    """

    HELP = "P or ENTER: pick up / H or *: hang up / Q: quit / 0-9: dial number / ?: Help"

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
        selection = readkey().lower()
        if selection in ["p", key.ENTER]:
            print("> Pick up phone")
            self.__controller.pick_up()
        elif selection in ["h", "*"]:
            print("> Hang up phone")
            self.__controller.hang_up()
        elif selection == "?":
            print(KeyboardReceiver.HELP)
        elif selection == "q":
            print("> Quit")
            self.__controller.quit()
            self.__do_quit = True
        elif str(selection).isdigit() and int(selection) > -1 and int(selection) <= 9:
            print(f"> Dial number {selection}")
            self.__controller.dial(int(selection))
        else:
            print("> Invalid command, type ? for help")


class DialPhoneReceiver(Thread):
    def __init__(self, controller: Controller):
        super().__init__()
        self.__controller: Controller = controller