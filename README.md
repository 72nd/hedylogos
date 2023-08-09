<div align="center" style="border-bottom: none">
  <h1>
    <img src="misc/logo.png" width="120"/>
    <br>
    hedylogos
  </h1>
  <h2>Build interactive audio scenarios using old dial phones or numeric keyboards</h2>
  <p><a href="README-de.md">Deutsche Version</a></p>
</div>

With Hedylogos it is possible to develop interactive audio formats. The navigation between the individual audio snippets is done by entering/dialling numbers, as is also known from telephone hotline menus. The software supports two different input modes for users: On the one hand, old dial telephones can be used for this purpose or the software can be controlled with the help of a [numeric keypad](https://de.wikipedia.org/wiki/Ziffernblock). The creation of the audio scenarios is relatively simple and can be done in a graphical editor as well as via text editor (more on this below). Many different scenarios are conceivable, both in terms of content and use, such as interactive stories or as a player in a museum, which plays different information depending on the button pressed.

## Funding

<img alt="Logo des Ministeriums für Wissenschaft Forschung und Kultur des Landes Brandenburg" src="misc/mwfk.png" width="200" style="align:left"/>

The project was funded by the Ministry of Science, Research and Culture of the State of Brandenburg.


## Requirements

Since Hedylogos is written in Python, the software runs on almost any computer. Please note that a [Raspberry Pi](https://en.wikipedia.org/wiki/Raspberry_Pi) is required to use a dial phone, as a normal laptop or computer does not have the necessary interfaces. If you don't want to go to the trouble, you can also use the dial pad mode on a conventional computer with a dedicated dial pad (like the one you can buy here at [Galaxus](https://www.galaxus.de/de/s1/product/logilink-id0120-nummernblock-kabellos-tastatur-12817754)) and a headset.

The installation of Python depends on the operating system used, but can be easily done with [this manual](https://python.land/installing-python). Since Python is widely used, any problems should be easily solved by using a search engine.

You should also be somewhat familiar with the command line. You can find an introduction here for [Windows](https://www.makeuseof.com/tag/a-beginners-guide-to-the-windows-command-line/) and here for [MacOS](https://www.makeuseof.com/tag/beginners-guide-mac-terminal/).


## Installation

First download the [latest version](https://github.com/72nd/hedylogos/releases/latest). Unpack the archive and navigate to the folder in the terminal.

```
python -m venv .venv
pip install .
```

Afterwards, Hedylogos can be started with the command `hedylogos`.


## Example

In the folder `example` there is an example scenario which covers the most important functions and features of Hedylogos. To start the example, the following command can be used.

```
hedylogos run-keyboard example/scenario.json
```


## Create a scenario

There are two ways to create a scenario for Hedylogos. On the one hand, there is a graphical editor, on the other hand, advanced users can also edit the JSON file directly. In the following, both ways will be explained.


### With the editor

To make the creation of the scenarios as easy as possible, there is a [graphical editor](https://72nd.github.io/hedylogos/editor/). You can find out more about how to use it in the editor itself. At the end, download the scenario as a file and save it on your computer.


### Manually using a JSON file

If you are familiar with JSON, you can also create the scenario directly in the editor. The meaning of the individual cases should be self-explanatory. Otherwise, it is worth taking a look at the [graphical editor](https://72nd.github.io/hedylogos/editor/) and the explanations of the fields contained therein. The software also offers the possibility to generate a template of the file.

```
hedylogos init scenario.json
```

## Validate the scenario file

When creating larger scenarios, small errors and mistakes can easily creep in. So that these do not only become apparent during use, validation routines are built into Hedylogos, which automatically detect the vast majority of problems when loading a scenario file. Only the check whether all audio files are present must be triggered manually with a command.

```
hedylogos check path/to/scenario.json
```

## Play the scenario

Hedylogos offers two different modes. One is with an old dial phone or with a keyboard or [numeric keypad](https://de.wikipedia.org/wiki/Ziffernblock).


### With the keyboard / numeric keypad

The execution is started with this command.

```
hedylogos run-keyboard path/to/scenario.json
```

The execution can be controlled with the following keys:

- `p` or `<ENTER>`: Starts the scenario. If the scenario is already being played back, the playback is interrupted and starts again at the starting point.
- `h`: Stops the playback of the scenario. Primarily simulates the moment when the telephone receiver is hung up and therefore has no real meaning for the keyboard mode.
- `q`: Stops the execution of the programme. To prevent visitors from triggering this action, it is recommended to use a dedicated numeric keypad.
- `0-9`: Dial a number.


### Using a dial telephone

To use the input of a dial phone, the library [RotaryPi](https://pypi.org/project/rotarypi/) is used. The prerequisite for this is the use of a [Raspberry Pi](https://www.raspberrypi.org/). More about the pin assignment can be found in the [RotaryPi documentation](https://rotarypi.readthedocs.io/en/latest/).