from model import Scenario

from pathlib import Path
from typing_extensions import Annotated

import typer

app = typer.Typer()

@app.command()
def run(
    path: Annotated[Path, typer.Argument(help="path to input file")]
):
    pass

if __name__ == "__main__":
    app()