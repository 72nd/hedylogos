from .model import Node, Scenario

from enum import Enum
import queue
from pathlib import Path
from threading import Thread
from typing import Optional

import simpleaudio


class _Action(str, Enum):
    """Actions which can be handled by the controller."""

    PICK_UP = "pick_up"
    HANG_UP = "hang_up"
    DIAL = "dial"
    QUIT = "quit"
    PLAYBACK_ENDED = "playback_ended"


class _Command:
    def __init__(self, action: _Action, number: Optional[int]=None):
        self.action = action
        self.number = number


class Player(Thread):
    def __init__(self, path: Path, queue: queue.Queue):
        super().__init__()
        self.__path = path
        self.__queue = queue
        self.__player: Optional[simpleaudio.PlayObject] = None
        self.__manual_stop: bool = False
    
    def run(self):
        wave = simpleaudio.WaveObject.from_wave_file(str(self.__path))
        self.__player = wave.play()
        self.__player.wait_done()
        if not self.__manual_stop:
            self.__queue.put(_Command(_Action.PLAYBACK_ENDED))
    
    def stop(self):
        if not self.__player:
            return
        self.__manual_stop = True
        self.__player.stop()


class Controller(Thread):
    def __init__(self, scenario: Scenario):
        super().__init__()
        self.__queue: queue.Queue = queue.Queue()
        self.__player: Optional[Player] = None
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
                self.__player = Player(self.__current_node.audio, self.__queue)
                self.__player.start()
            elif command.action is _Action.HANG_UP:
                self.__current_node = None
                if self.__player:
                    self.__player.stop()
            elif command.action is _Action.DIAL:
                if not self.__current_node:
                    continue
                target = self.__current_node.next_node_id_by_number(command.number)
                if not target:
                    print("invalid number, TODO: implement handling")
                    continue
                if self.__player:
                    self.__player.stop()
                self.__current_node = self.__scenario.get_nodes_dict()[target]
                self.__player = Player(self.__current_node.audio, self.__queue)
                self.__player.start()
            elif command.action is _Action.PLAYBACK_ENDED:
                print("ended by itself")
            elif command.action is _Action.QUIT:
                if self.__player:
                    self.__player.stop()
                break