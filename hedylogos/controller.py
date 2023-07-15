from .model import Node, Scenario
from .player import Player

from enum import Enum
import queue
from threading import Thread
from typing import Optional


class _Action(str, Enum):
    """Actions which can be handled by the controller."""

    PICK_UP = "pick_up"
    HANG_UP = "hang_up"
    DIAL = "dial"
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
        self.__player.start()
        self.__scenario = scenario
        self.__current_node: Optional[Node] = None
    
    def pick_up(self):
        """Someone picked up the phone."""
        self.__queue.put(_Command(_Action.PICK_UP))

    def hang_up(self):
        """The phone was hung up therefore ending the execution and resetting."""
        self.__queue.put(_Command(_Action.HANG_UP))

    def dial(self, number: int):
        """A number input occurred."""
        self.__queue.put(_Command(_Action.DIAL, number))
    
    def quit(self):
        """Quit execution."""
        self.__queue.put(_Command(_Action.QUIT))

    def run(self):
        while True:
            command = self.__queue.get()
            if command.action is _Action.PICK_UP:
                self.__current_node = self.__scenario.start()
                self.__player.play(self.__current_node.audio)
            elif command.action is _Action.HANG_UP:
                self.__current_node = None
                self.__player.stop()
            elif command.action is _Action.DIAL:
                if not self.__current_node:
                    continue
                target = self.__current_node.next_node_id_by_number(command.number)
                if not target:
                    print("invalid number, TODO: implement handling")
                    continue
                self.__player.stop()
                self.__current_node = self.__scenario.get_nodes_dict()[target]
                self.__player.play(self.__current_node.audio)
            elif command.action is _Action.QUIT:
                self.__player.quit()
                break