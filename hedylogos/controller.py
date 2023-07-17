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
    """
    The controller handles the playback of the audio files and reacts to events
    with the phone.
    """

    PLAY_NORMAL_MESSAGE = 2
    """
    How many times should the normal normal invalid number audio be played before
    the fun message is played (if available). THus if PLAY_FUN_MESSAGE is set to
    3 the »normal message« will be played the first three occurrence and on the
    fourth time the fun message will be played.
    """

    def __init__(self, scenario: Scenario, scenario_path: Path):
        super().__init__()
        self.__scenario = scenario
        self.__scenario_location = scenario_path.resolve().parent
        self.__queue: queue.Queue = queue.Queue()
        self.__player: Optional[Player] = None
        self.__current_node: Optional[Node] = None
        self.__plays_invalid_audio: bool = False
        self.__played_normal_invalid_audio: int = 0
    
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
                self.__on_pick_up()
            elif command.action is _Action.HANG_UP:
                self.__on_hang_up()
            elif command.action is _Action.DIAL:
                self.__on_dial(command)
            elif command.action is _Action.PLAYBACK_ENDED:
                self.__on_playback_ended()
            elif command.action is _Action.QUIT:
                self.__on_quit()
                break
        
    def check_paths(self):
        errors_present: bool = False
        paths: list[str] = [
            self.__scenario.invalid_number_audio,
            self.__scenario.internal_error_audio,
            self.__scenario.end_call_audio,
        ]
        if self.__scenario.invalid_number_fun_audio:
            paths.append(self.__scenario.invalid_number_fun_audio)
        paths.extend([node.audio for node in self.__scenario.nodes.root])
        for path in paths:
            absolute_path = self.__resolve_path(path)
            if not absolute_path.exists():
                print(f"{absolute_path} not found")
                errors_present = True
        if not errors_present:
            print("All paths are valid.")
    
    def __on_pick_up(self):
        if self.__player:
            self.__player.stop()
        self.__current_node = self.__scenario.start()
        self.__start_playback(self.__current_node.audio)

    def __on_hang_up(self):
        self.__current_node = None
        if self.__player:
            self.__player.stop()

    def __on_dial(self, command: _Command):
        if not self.__current_node:
            return
        if not command.number:
            raise RuntimeError("__on_dial called without a number")
        target = self.__current_node.next_node_id_by_number(command.number)
        if not target:
            self.__on_invalid_number()
            return
        if self.__player:
            self.__player.stop()
        self.__current_node = self.__scenario.node_by_id(target)
        self.__start_playback(self.__current_node.audio)
    
    def __on_invalid_number(self):
        if self.__player:
            self.__player.stop()
        path = self.__scenario.invalid_number_audio
        # See documentation of `Controller.PLAY_NORMAL_MESSAGE` for more
        # Information about this. Why? Because it's fun, that's why.
        if self.__scenario.invalid_number_fun_audio:
            if self.__played_normal_invalid_audio == Controller.PLAY_NORMAL_MESSAGE:
                path = self.__scenario.invalid_number_fun_audio
                self.__played_normal_invalid_audio = 0
            else:
                self.__played_normal_invalid_audio += 1
        self.__start_playback(path)
        self.__plays_invalid_audio = True

    def __on_playback_ended(self):
        if not self.__current_node:
            # Block only here vor more clarity, happens when a scenario ended.
            pass
        elif self.__plays_invalid_audio:
            # Invalid audio playback ended, replay the current node.
            self.__start_playback(self.__current_node.audio)
            self.__plays_invalid_audio = False
        elif not self.__current_node.links:
            # The current node has no links defined so the scenario execution ends.
            self.__current_node = None
            self.__start_playback(self.__scenario.end_call_audio)
        elif self.__current_node.has_unnumbered_links():
            selected_link = self.__current_node.random_link()
            if not selected_link:
                raise RuntimeError("tried to get an random link on a node with no links")
            self.__current_node = self.__scenario.node_by_id(selected_link.target)
            self.__start_playback(self.__current_node.audio)

    def __on_quit(self):
        if self.__player:
            self.__player.stop()
    
    def __start_playback(self, path: str):
        self.__player = Player(
            self.__resolve_path(path),
            self.__queue,
        )
        self.__player.start()
    
    def __resolve_path(self, path: str):
        """Used to resolve paths relative to the scenario file."""
        return (self.__scenario_location / Path(path)).resolve()