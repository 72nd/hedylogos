from enum import Enum
import queue
from pathlib import Path
from threading import Thread
import time
from typing import Optional, Tuple
import wave

import pyaudio


class _PlayerAction(str, Enum):
    """Commands which can be sent to the player thread."""
    PLAY = "play"
    STOP = "stop"
    QUIT = "quit"


class _PlayerCommand:
    def __init__(self, action: _PlayerAction, path: Optional[Path]):
        self.action = action
        self.path = path


class Player(Thread):
    def __init__(self):
        super().__init__() 
        self.__queue = queue.Queue()
        self.__pyaudio: pyaudio.PyAudio = pyaudio.PyAudio()
        self.__wave: Optional[wave.Wave_read] = None
        self.__stream: Optional[pyaudio.Stream] = None

    def play(self, path: Path):
        self.__queue.put(_PlayerCommand(_PlayerAction.PLAY, path))
    
    def stop(self):
        self.__queue.put(_PlayerCommand(_PlayerAction.STOP, None))
    
    def quit(self):
        self.__queue.put(_PlayerCommand(_PlayerAction.QUIT, None))
    
    def run(self):
        while True:
            try:
                command = self.__queue.get_nowait()
            except queue.Empty:
                command = None
            if self.__stream and self.__stream.is_active():
                time.sleep(0.1)
            if not command:
                continue
            if command.action is _PlayerAction.PLAY:
                self.__play(command.path)
            elif command.action is _PlayerAction.STOP:
                self.__stop()
            elif command.action is _PlayerAction.QUIT:
                print(quit)
                self.__stop()
                self.__pyaudio.terminate()
                break
    
    def __play(self, path: Path):
        if self.__wave or self.__stream:
            raise RuntimeError("stop running audio before starting another one")
        self.__wave = wave.open(str(path), "rb")
        self.__stream = self.__pyaudio.open(
            format=self.__pyaudio.get_format_from_width(self.__wave.getsampwidth()),
            channels=self.__wave.getnchannels(),
            rate=self.__wave.getframerate(),
            output=True,
            stream_callback=self.__callback
        )
        self.__stream.start_stream()

    def __callback(self, in_data, frame_count, time_info, status):
        if not self.__wave:
            raise RuntimeError("pyaudio callback called while no PyAudio instance is present")
        return (
            self.__wave.readframes(frame_count),
            pyaudio.paContinue,
        )
    
    def __stop(self):
        print("stop")
        if not self.__stream or not self.__wave:
            return
        self.__stream.stop_stream()
        self.__stream.close()
        self.__wave.close()
        self.__stream = None
        self.__wave = None