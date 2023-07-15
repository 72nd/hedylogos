from collections import Counter
import json
from pathlib import Path

from pydantic import BaseModel, Field, RootModel, field_validator
from typing import Optional


def format_str_list(data: list[str]):
    return f"\'{', '.join(data)}\'"


class Link(BaseModel):
    """Used to define next nodes in a scenario."""

    target: str = Field(examples=["node_id"])
    """Id of the node which the link should point to."""
    number: Optional[int] = Field(ge=0, le=9, examples=[0])
    """A number between 0 and 9 for the user to be able to choose the link."""


class Node(BaseModel):
    """
    A node represents a state within the graph and contains a audio.
    """

    id: str = Field(examples=["node_id"])
    """
    Unique identifier of the node. Used to refer to the node in other parts of the scenario.
    """
    name: str = Field(examples=["A node"])
    """Name of the node for debugging purposes."""
    content: Optional[str]
    """Optional the textual content of the audio in the node."""
    audio: Path
    """Path to the audio file."""
    links: Optional[list[Link]]
    """Links to other nodes next in the story/scenario line."""

    @classmethod
    def start_example(cls) -> "Node":
        return cls(
            id="start",
            name="Start Node",
            content=None,
            audio=Path("start.mp3"),
            links=[Link(
                target="start",
                number=1,
            )]
        )
    
    @classmethod
    def additional_example(cls) -> "Node":
        return cls(
            id="another",
            name="Another node",
            content="Hello this is the content text",
            audio=Path("another.mp3"),
            links=None,
        )

    def next_node_id_by_number(self, number: int) -> Optional[str]:
        """Returns the id of a linked node if the number is valid."""
        if not self.links:
            return None 
        for link in self.links:
            if link.number == number:
                return link.target
        return None

class Nodes(RootModel[list[Node]]):
    """
    A collection of nodes. Has it's own class to write a validator which
    makes sure all `Node.id` are unique.
    """
    root: list[Node]

    @classmethod
    def init_example(cls) -> "Nodes":
        """Returns a instance of the model with some initial data."""
        return cls(
            root=[Node.start_example()]
        )

    @field_validator("root")
    def check_unique_ids(cls, v):
        ids = [node.id for node in v]
        counter = Counter(ids)
        duplicates = [id for id, count in counter.items() if count > 1]
        if len(duplicates) > 0:
            raise ValueError(f"node id(s) {format_str_list(duplicates)} not unique")
        return v

    @field_validator("root")
    def check_valid_link_targets(cls, v):
        ids = [node.id for node in v]
        for node in v:
            if node.links is None:
                continue
            invalid = [link.target for link in node.links if link.target not in ids]
            if len(invalid) != 0:
                raise ValueError(f"node '{node.id}' has invalid link target(s) {format_str_list(invalid)}")
        return v
    
    def as_dict(self) -> dict[str, Node]:
        """Returns a dict of the Nodes with id as their keys."""
        rsl: dict[str, Node] = {}
        for node in self.root:
            rsl[node.id] = node
        return rsl
    

class Scenario(BaseModel):
    """
    The scenario represents the base data structure of an story scenario.
    It contains all metadata and steps within the story.
    """

    name: str = Field(examples=["The story of the hotline"])
    """The name of the scenario."""
    description: Optional[str] = Field(examples=["Some information about the scenario"])
    """Gives some information about the scenario defined by the graph."""
    authors: list[str]
    """A list of the names."""
    nodes: Nodes
    """All nodes of the scenario."""
    start_node: str
    """The id of the node which the scenario should start with."""
    nodes_dict: Optional[dict[str, Node]] = Field(exclude=True, default=None)

    @classmethod
    def from_json(cls, path: Path) -> "Scenario":
        """Reads a `Scenario` instance from a JSON file."""
        with open(path, "r") as f:
            return cls.model_validate_json(f.read())
    
    @classmethod
    def init_example(cls) -> "Scenario":
        """Returns a instance of the model with some initial data."""
        return cls(
            name="A interactive Scenario",
            description="This is an almost empty scenario file",
            authors=["Max Mustermann"],
            nodes=Nodes.init_example(),
            start_node="start"
        )
    
    @classmethod
    def to_json_schema(cls, path: Path):
        """
        Saves the the JSON schema of Scenario to a file so it can be used in
        a graphical JSON editor.
        """
        with open(path, "w") as f:
            f.write(json.dumps(cls.model_json_schema()))

    def to_json(self, path: Path):
        """Write the scenario instance to a JSON file."""
        with open(path, "w") as f:
            f.write(self.model_dump_json())
    
    def start(self) -> Node:
        if self.start_node not in self.get_nodes_dict():
            raise KeyError(f"no Node for start_node '{self.start_node} found")
        return self.get_nodes_dict()[self.start_node]

    def get_nodes_dict(self) -> dict[str, Node]:
        if not self.nodes_dict:
            self.nodes_dict = self.nodes.as_dict()
        return self.nodes_dict