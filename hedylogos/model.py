from collections import Counter
from dataclasses import field
import json
from pathlib import Path
import random

from pydantic import BaseModel, Field, RootModel, field_validator
from typing import Optional


def format_str_list(data: list[str]):
    return f"\'{', '.join(data)}\'"


class Link(BaseModel):
    """Used to define next nodes in a scenario."""

    target: str = Field(
        examples=["node_id"],
        description="Id of the node which the link should point to.",
        min_length=1
    )
    number: Optional[int] = Field(
        ge=0,
        le=9,
        examples=[0],
        description="A number between 0 and 9 for the user to be able to choose the link. If set to null a random link will be chosen."
    )


class Node(BaseModel):
    """
    A node represents a state within the graph and contains a audio.
    """

    id: str = Field(
        examples=["node_id"],
        description="Unique identifier of the node. Used to refer to the node in other parts of the scenario.",
        min_length=1
    )
    name: str = Field(
        examples=["A node"],
        description="Name of the node for debugging purposes.",
        min_length=1
    )
    content: Optional[str] = Field(
        description="Optional the textual content of the audio in the node.",
        json_schema_extra={
            "format": "textarea",
        }
    )
    audio: str = Field(
        description="Path to the audio file. Can be relative to the location of the scenario file.",
        min_length=1
    )
    links: Optional[list[Link]] = Field(
        description="Links to other nodes next in the story/scenario line. If set to None the scenario will stop at this point."
    )

    @classmethod
    def start_example(cls) -> "Node":
        return cls(
            id="start",
            name="Start Node",
            content=None,
            audio="start.mp3",
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
            audio="another.mp3",
            links=None,
        )

    @field_validator("links")
    def check_only_one_link_type(cls, v):
        if not v or len(v) <= 1:
            return v
        is_none = False
        if v[0].number is None:
            is_none = True
        for link in v:
            if link.number is None and not is_none:
                raise ValueError("mixed use of numbers and None in links")
        return v

    def next_node_id_by_number(self, number: int) -> Optional[str]:
        """Returns the id of a linked node if the number is valid."""
        if not self.links:
            return None 
        for link in self.links:
            if link.number == number:
                return link.target
        return None
    
    def random_link(self) -> Optional[Link]:
        """Selects randomly one of the links of the node."""
        if not self.links:
            return None
        return random.choice(self.links)
    
    def has_unnumbered_links(self) -> bool:
        """Returns whether links have numbers assigned to them."""
        if not self.links:
            return False
        return len([link for link in self.links if not link.number]) > 0


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

    name: str = Field(
        examples=["The story of the hotline"],
        description="The name of the scenario.",
        min_length=1,
    )
    description: Optional[str] = Field(
        examples=["Some information about the scenario"],
        description="Gives some information about the scenario defined by the graph."
    )
    authors: Optional[list[str]] = Field(
        description="A list of the names.",
        default=None,
    )
    nodes: Nodes = Field(
        description="All nodes of the scenario."
    )
    start_node: str = Field(
        description="Id of the node the execution should start.",
        min_length=1,
    )
    invalid_number_audio: str = Field(
        description="Audio played when the user dials an invalid number.",
        min_length=1,
    )
    invalid_number_fun_audio: Optional[str] = Field(
        description="Optional fun audio played some times when the user dials an invalid number. ",
        min_length=1,
    )
    internal_error_audio: str = Field(
        description="Audio played when an internal error occurred.",
        min_length=1,
    )
    end_call_audio: str = Field(
        description="Audio played when scenario ends.",
        min_length=1,
    )
    nodes_dict: Optional[dict[str, Node]] = Field(exclude=True, default=None)
    """The id of the node which the scenario should start with."""

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
            start_node="start",
            invalid_number_audio="invalid-number.wav",
            invalid_number_fun_audio=None,
            internal_error_audio="error.wav",
            end_call_audio="end-call.wav",
            nodes=Nodes.init_example(),
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
        """Returns the start node."""
        if self.start_node not in self.get_nodes_dict():
            raise KeyError(f"no Node for start_node '{self.start_node} found")
        return self.get_nodes_dict()[self.start_node]

    def get_nodes_dict(self) -> dict[str, Node]:
        """Returns a dict of all nodes with their id as key."""
        if not self.nodes_dict:
            self.nodes_dict = self.nodes.as_dict()
        return self.nodes_dict
    
    def node_by_id(self, id: str) -> Node:
        try:
            return self.get_nodes_dict()[id]
        except KeyError:
            raise KeyError(f"there is no node with '{id}' in the scenario")