package model

import (
	"github.com/freddy33/graphml"
)

type Node struct {
	ID string
}

func NewNode(obj graphml.Node, keys Keys) (*Node, error) {
	return &Node{
		ID: obj.ID,
	}, nil
}

type Nodes []Node

func NewNodes(graph graphml.Graph, keys Keys) (*Nodes, error) {
	var rsl Nodes
	for _, node := range graph.Nodes {
		item, err := NewNode(node, keys)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *item)
	}
	return &rsl, nil
}
