from model import Scenario

from pathlib import Path
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
    pass


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


if __name__ == "__main__":
    app()