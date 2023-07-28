from .controller import Controller
from .model import Scenario
from .receiver import KeyboardReceiver

from pathlib import Path

from typing_extensions import Annotated

import typer

app = typer.Typer()


@app.command()
def check(
    path: Annotated[Path, typer.Argument(help="output path")]
):
    """
    Checks if all audio files in the scenario can be found
    """
    scenario = Scenario.from_json(path)
    controller = Controller(scenario, path)
    controller.check_paths()


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
    controller = Controller(scenario, path)
    receiver = KeyboardReceiver(controller)
    receiver.run()


@app.command()
def run_phone(
    path: Annotated[Path, typer.Argument(help="path to scenario file")]
):
    """
    Runs the scenario using input form a dial phone using rotarypi.
    """
    scenario = Scenario.from_json(path)
    controller = Controller(scenario, path)
    if Path("/etc/rpi-issue").exists():
        from .receiver import DialPhoneReceiver
        receiver = DialPhoneReceiver(controller)
        receiver.run()
    else:
        raise NotImplementedError("Reading input from a rotary phone is only implemented for Raspberry Pi's (where the rp.GPIO library is available)")


@app.command()
def schema(
    path: Annotated[Path, typer.Argument(help="output path")]
):
    """
    Writes the JSON Schema of a Scenario file to disk.
    """
    Scenario.to_json_schema(path)


if __name__ == "__main__":
    app()