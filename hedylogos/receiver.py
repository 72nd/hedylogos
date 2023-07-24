import queue
from .controller import Controller

import logging
from queue import Queue
from threading import Thread

from readchar import readkey, key
from rotarypi import DialEvent, DialConfiguration, DialPinout, EventType, HandsetState, RotaryReader


class KeyboardReceiver(Thread):
    """
    Provides a simple shell for entering the different commands to the system. This receiver
    is mainly used for debugging. This mode can also be used with an numeric-keyboard to provide
    a simple user facing control panel without the need of setting up a real old dial phone.
    """

    HELP = "P or ENTER: pick up / H: hang up / Q: quit / 0-9: dial number / ?: Help"

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
        elif selection == "h":
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
        self.__queue: Queue[DialEvent] = Queue()
        self.__reader: RotaryReader = RotaryReader(
            self.__queue,
            DialPinout(),
            DialConfiguration(loglevel=logging.DEBUG),
        )
    
    def run(self):
        self.__reader.start()
        self.__controller.start()
        while True:
            self.__on_event(self.__queue.get())
            
    def __on_event(self, event: DialEvent):
        if event.type == EventType.DIAL_EVENT:
            if not isinstance(event.data, int):
                print("Logic error: got DIAL_EVENT with HandsetData as data")
                return
            print(f"> Dial number {event.data}")
            self.__controller.dial(event.data)
        if event.type == EventType.HANDSET_EVENT:
            if event.data is HandsetState.PICKED_UP:
                print("> Pick up phone")
                self.__controller.pick_up()
            elif event.data is HandsetState.HUNG_UP:
                print("> Hang up phone")
                self.__controller.hang_up()