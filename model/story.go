// The model package contains the definition of the fundamental data types
// used to represent the structure of the story represented as a graph.
// This package also handles the import and conversion from the more basic
// result of an [github.com/hedylogos/graphml] import.l
package model

import (
	"github.com/72nd/hedylogos/graphml"
)

// Story contains a whole audio scenario containing the metadata and the
// start nodes.
type Story struct {
	// Optional version number for a story. This is used to distinguish
	// between different versions of a stroy.
	Version string
	// The optional name(s) of the author(s) of the story.
	Author string
	// Optional description of the content of the story.
	Description string
	// Languages supported by the story.
	Languages Languages
	// The start nodes of the story.
	Starts []*Node
	// The collection of all nodes in the story. Should be only used
	// internally.
	Nodes Nodes
	// All keys of the story.
	Keys Keys
}

// Takes a [github.com/72nd/hedylogos/graphml.Document] validates the
// data and returns a new [Story] instance based on it.
func NewStory(doc graphml.Document) (*Story, error) {
	keys, err := NewKeys(doc.Keys)
	if err != nil {
		return nil, err
	}
	sto, err := NewStorage(doc.Graph.Data, *keys)
	if err != nil {
		return nil, err
	}
	version, err := ValueByName[string](*sto, "Version")
	if err != nil {
		return nil, err
	}
	author, err := ValueByName[string](*sto, "Author")
	if err != nil {
		return nil, err
	}
	desc, err := ValueByName[string](*sto, "Description")
	if err != nil {
		return nil, err
	}
	langs, err := ValueByName[Languages](*sto, "Languages")
	if err != nil {
		return nil, err
	}
	nodes, err := NewNodes(doc.Graph.Nodes, *keys)
	if err != nil {
		return nil, err
	}
	return &Story{
		Version:     *version,
		Author:      *author,
		Description: *desc,
		Keys:        *keys,
		Languages:   *langs,
		Nodes:       *nodes,
	}, nil
}
