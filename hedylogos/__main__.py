from .model import Scenario
from .player import PlayerAction, Player

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


@app.command()
def run_mqtt(
    path: Annotated[Path, typer.Argument(help="path to scenario file")]
):
    """
    Runs the scenario using the input form the dial via the MQTT server.
    """
    pass


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
    player.send_command(PlayerAction.PLAY, path)
    time.sleep(10)
    player.send_command(PlayerAction.STOP)
    time.sleep(1)
    player.send_command(PlayerAction.QUIT)


if __name__ == "__main__":
    app()