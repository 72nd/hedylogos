import click

from model import Scenario

from pathlib import Path

@click.command()
@click.option("--input", "-i", "input_file", help="path to input file", required=True)
def run(input_file: Path):
    scenario = Scenario(input_file)
    print(scenario['title'])


if __name__ == "__main__":
    run()