from pathlib import Path
from pydantic import BaseModel, Field
from typing import Optional

from exception import NoDiGraph

from pathlib import Path


class Link(BaseModel):
    """Used to define next nodes in a scenario."""

    class Config:
        populate_by_name = True

    target: str
    """Id of the node which the link should point to."""
    number: int
    """A number between 0 and 9 for the user to be able to choose the link."""


class Node(BaseModel):
    """
    A node represents a state within the graph and contains a audio.
    """

    id: str = Field()
    """
    Unique identifier of the node. Used to refer to the node in other parts of the scenario.
    """
    name: str
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


class Scenario(BaseModel):
    """
    The scenario represents the base data structure of an story scenario.
    It contains all metadata and steps within the story.
    """

    name: str
    """The name of the scenario."""
    description: Optional[str]
    """Gives some information about the scenario defined by the graph."""
    authors: list[str]
    """A list of the names."""
    nodes: list[Node]
    """All nodes of the scenario."""
    start_node: str
    """The id of the node which the scenario should start with."""

    @classmethod
    def from_json(cls, path: Path) -> "Scenario":
        """Reads a `Scenario` instance from a file."""
        with open(path, "r") as f:
            return cls.model_validate_json(f.read())
    
    @classmethod
    def init_example(cls) -> "Scenario":
        """Returns a instance of the model with some initial data."""
        return cls(
            name="A interactive Scenario",
            description="This is an almost empty scenario file",
            authors=["Max Mustermann"],
            nodes=[Node.start_example()],
            start_node="start"
        )
    
    def to_json(self, path: Path):
        with open(path, "r") as f:
            f.write(self.model_dump_json(by_alias=True))