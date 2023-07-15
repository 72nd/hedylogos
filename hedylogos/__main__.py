from .controller import Controller
from .model import Scenario
from .player import Player
from .receiver import KeyboardReceiver, MqttReceiver

from pathlib import Path
import time
from typing_extensions import Annotated

import typer

app = typer.Typer()

@app.command()
def init(
    path: Annotated[Path, typer.Argument(help="output path")]
):
    """
    Writes a new JSON Scenario files with some example values to the disk.
    """
    scenario = Scenario.init_example()
    scenario.to_json(path)


@app.command()
def run_keyboard(
    path: Annotated[Path, typer.Argument(help="path to scenario file")]
):
    """
    Runs the scenario using the input form the keyboard.
    """
    scenario = Scenario.from_json(path)
    controller = Controller(scenario)
    receiver = KeyboardReceiver(controller)
    receiver.run()


@app.command()
def run_mqtt(
    path: Annotated[Path, typer.Argument(help="path to scenario file")]
):
    """
    Runs the scenario using the input form the dial via the MQTT server.
    """
    scenario = Scenario.from_json(path)
    controller = Controller(scenario)
    receiver = MqttReceiver(controller)
    receiver.run()


@app.command()
def schema(
    path: Annotated[Path, typer.Argument(help="output path")]
):
    """
    Writes the JSON Schema of a Scenario file to disk.
    """
    Scenario.to_json_schema(path)


@app.command()
def test(
    path: Annotated[Path, typer.Argument(help="audio file")]
):
    """It's a test."""
    player = Player()
    player.start()
    player.play(path)
    time.sleep(4)
    player.stop()
    player.play(Path("share/error.wav"))
    time.sleep(5)
    player.quit()


if __name__ == "__main__":
    app()