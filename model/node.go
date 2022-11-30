package model

import (
	"strings"

	"github.com/72nd/hedylogos/enum"
	"github.com/72nd/hedylogos/graphml"
)

// Defines the different node types and thus the behavior of the node.
type NodeType string

const (
	StartNodeType   = NodeType("start")
	EndNodeType     = NodeType("end")
	LangNodeType    = NodeType("language")
	MenuNodeType    = NodeType("menu")
	ExecuteNodeType = NodeType("execute")
	IfNodeType      = NodeType("if")
)

func (t NodeType) Instances() []NodeType {
	return []NodeType{StartNodeType, EndNodeType, LangNodeType, MenuNodeType, ExecuteNodeType, IfNodeType}
}

// The different shapes a node can have in it's visual representation
// within yEd. This information is used to enforce the correct usage
// of node shapes within the input file.
type NodeShape string

const (
	EllipseShape        = NodeShape("ellipse")
	StarShape           = NodeShape("star5")
	RoundRectangleShape = NodeShape("roundrectangle")
	HexagonShape        = NodeShape("hexagon")
)

func (s NodeShape) Instances() []NodeShape {
	return []NodeShape{EllipseShape, StarShape, RoundRectangleShape, HexagonShape}
}

// Implemented by all elements which can be child of node's output element.
// Used to define the output of any given node. This can be something like
// audio or text-to-speech (TTS) output.
type Output interface {
	Execute()
}

// Collection of [Node]s.
type Nodes []Node

// Takes a slice of `graphml.Node`s and returns a new [Nodes] instance.
func NewNodes(nodes []graphml.Node, keys Keys) (*Nodes, error) {
	var rsl Nodes
	for _, node := range nodes {
		nd, err := NewNode(node, keys)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *nd)
	}
	return &rsl, nil
}

// A node represent one logic step within the story.
type Node struct {
	// Unique identifier of the node. Given by the GraphML.
	ID string
	// The name is based on the label content provided by yEd.
	Name string
	// The type of the node. This defines the behavior of the node
	// when called.
	Type NodeType
	// Shape used for the visual representation of the node within yEd.
	Shape NodeShape
}

// Returns a new node based on the graphml version.
func NewNode(node graphml.Node, keys Keys) (*Node, error) {
	sto, err := NewStorage(node.Data, keys)
	if err != nil {
		return nil, err
	}
	tpEle, err := ValueByName[string](*sto, "Type")
	if err != nil {
		return nil, err
	}
	tp, err := enum.EnumByValue[NodeType](StartNodeType, NodeType(strings.ToLower(*tpEle)))
	if err != nil {
		return nil, err
	}
	shapeEle, err := ValueByName[graphml.ShapeData](*sto, "yft.nodegraphics")
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(shapeEle.Label)
	shape, err := enum.EnumByValue[NodeShape](EllipseShape, NodeShape(shapeEle.Shape.Type))
	if err != nil {
		return nil, err
	}
	return &Node{
		ID:    node.ID,
		Name:  name,
		Type:  *tp,
		Shape: *shape,
	}, nil
}
