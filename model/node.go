package model

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/72nd/nohelpline/enum"
	"github.com/freddy33/graphml"
)

// Collection of the different node types
type NodeType string

const (
	StartNodeType   = NodeType("start")
	EndNodeType     = NodeType("end")
	LangNodeType    = NodeType("language")
	MenuNodeType    = NodeType("menu")
	ExecuteNodeType = NodeType("execute")
	IfNodeType      = NodeType("if")
)

func (n NodeType) Instances() []NodeType {
	return []NodeType{StartNodeType, EndNodeType, LangNodeType, MenuNodeType, ExecuteNodeType, IfNodeType}
}

// Implemented by all elements which can be child of node's output element.
// Used to define the output of any given node. This can be something like
// audio or text-to-speech (TTS) output.
type Output interface {
	Execute()
}

// Basic attributes which each node contains.
type NodeBase struct {
	ID   string
	Name string
	Type NodeType
}

func NewNodeBase(obj graphml.Node, keys Keys, ntype NodeType) (*NodeBase, error) {
	name, err := yedNodeName(obj, keys)
	if err != nil {
		return nil, err
	}
	return &NodeBase{
		ID:   obj.ExtObject.ID,
		Name: name,
		Type: ntype,
	}, nil
}

func yedNodeName(obj graphml.Node, keys Keys) (string, error) {
	ng, err := keys.DataByName(obj.ExtObject, string(NodeGraphicsKey))
	if err != nil {
		return "", err
	}
	insideNodeLabel := false
	for _, ele := range ng {
		item, ok := ele.(xml.StartElement)
		if ok && item.Name.Local == "NodeLabel" {
			insideNodeLabel = true
			continue
		}
		value, ok := ele.(string)
		if ok && insideNodeLabel {
			return value, nil
		}
	}
	return "<UNDEFINED>", fmt.Errorf("node with id %s has no label", obj.ID)
}

// A node represents a stage in the control flow of the stroy. There are
// different types of nodes which represents different aspects of the
// control flow.
type Node interface {
	ID() string
	Type() NodeType
}

// Parses the nodes of the graphml file and casts them into the correct
// node type.
func NewNode(obj graphml.Node, keys Keys) (Node, error) {
	nType, err := keys.DataByName(obj.ExtObject, string(TypeKey))
	if err != nil {
		return nil, err
	}
	snType, ok := nType[0].(string)
	if !ok {
		return nil, fmt.Errorf("type data in node has to be a string")
	}
	snType = strings.ToLower(snType)
	tnType, err := enum.EnumByValue[NodeType](StartNodeType, NodeType(snType))
	if err != nil {
		return nil, err
	}
	switch *tnType {
	case StartNodeType:
		return NewStartNode(obj, keys, *tnType)
	case EndNodeType:
		return NewEndNode(obj, keys, *tnType)
	default:
		return newNotImplementedNode(obj, keys, *tnType)
	}
	return nil, fmt.Errorf("NewNode isn't currently implemented to handle NodeType %s", *tnType)
}

// Entry point for an story scenario. Other than that this node doesn't
// provides any additional functionality.
type StartNode struct {
	NodeBase
}

// Parses a given a graphml node and returns it as a start node. The method
// expects a node with the property `Type` set to `start`.
func NewStartNode(obj graphml.Node, keys Keys, ntype NodeType) (*StartNode, error) {
	base, err := NewNodeBase(obj, keys, ntype)
	if err != nil {
		return nil, err
	}
	return &StartNode{
		NodeBase: *base,
	}, nil
}

func (s StartNode) ID() string {
	return s.NodeBase.ID
}

func (s StartNode) Type() NodeType {
	return s.NodeBase.Type
}

// Entry point for an story scenario. Other than that this node doesn't
// provides any additional functionality.
type EndNode struct {
	NodeBase
}

// Parses a given a graphml node and returns it as a end node. The method
// expects a node with the property `Type` set to `end`.
func NewEndNode(obj graphml.Node, keys Keys, ntype NodeType) (*EndNode, error) {
	base, err := NewNodeBase(obj, keys, ntype)
	if err != nil {
		return nil, err
	}
	return &EndNode{
		NodeBase: *base,
	}, nil
}

func (e EndNode) ID() string {
	return e.NodeBase.ID
}

func (e EndNode) Type() NodeType {
	return e.NodeBase.Type
}

// Used for development only. Will be removed.
type notImplementedNode struct {
	NodeBase
}

func newNotImplementedNode(obj graphml.Node, keys Keys, ntype NodeType) (*notImplementedNode, error) {
	base, err := NewNodeBase(obj, keys, ntype)
	if err != nil {
		return nil, err
	}
	return &notImplementedNode{
		NodeBase: *base,
	}, nil
}

func (n notImplementedNode) ID() string {
	return n.NodeBase.ID
}

func (n notImplementedNode) Type() NodeType {
	return n.NodeBase.Type
}

type Nodes []Node

func NewNodes(graph graphml.Graph, keys Keys) (*Nodes, error) {
	var rsl Nodes
	for _, node := range graph.Nodes {
		item, err := NewNode(node, keys)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, item)
	}
	return &rsl, nil
}
