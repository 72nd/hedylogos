from pathlib import Path
class NoDiGraph(Exception):
    """Exception raised when a given file doesn't contain directed graph."""

    def __init__(self):
        pass
