import networkx as nx

from exception import NoDiGraph

from pathlib import Path


class Scenario(nx.DiGraph):
    def __init__(self, path: Path):
        self.source = nx.read_graphml(path)
