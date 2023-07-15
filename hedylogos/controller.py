from .model import Scenario
from .player import Player

from enum import Enum
import queue
from threading import Thread
from typing import Optional


class _Action(str, Enum):
    """Actions which can be handled by the controller."""

    PICK_UP = "pick_up"
    HANG_UP = "hang_up"
    NUMBER = "number"
    QUIT = "quit"


class _Command:
    def __init__(self, action: _Action, number: Optional[int]=None):
        self.action = action
        self.number = number


class Controller(Thread):
    def __init__(self, scenario: Scenario):
        super().__init__()
        self.__queue = queue.Queue()
        self.__player = Player()
        self.__scenario = scenario
    
    def pick_up(self):
        """Someone picked up the phone."""
        self.__queue.put(_Command(_Action.PICK_UP))

    def hang_up(self):
        """The phone was hung up therefore ending the execution and resetting."""
        self.__queue.put(_Command(_Action.HANG_UP))

    def number(self, number: int):
        """A number input occurred."""
        self.__queue.put(_Command(_Action.NUMBER, number))
    
    def quit(self):
        """Quit execution."""
        self.__queue.put(_Command(_Action.QUIT))


    def run(self):
        while True:
            command = self.__queue.get()
            match command.action:
                case _Action.PICK_UP:
                    # TODO Start with scenario
                    pass
                case _Action.HANG_UP:
                    self.__player.stop()
                case _Action.NUMBER:
                    # TODO: Play the correct file
                    pass
                case _Action.QUIT:
                    self.__player.quit()