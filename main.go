package main

import (
	"fmt"
	"log"
	"os"

	"github.com/72nd/nohelpline/enum"
	"github.com/freddy33/graphml"
)

type KeyName string

const (
	LanguagesKey   = KeyName("Languages")
	DescriptionKey = KeyName("Description")
	AuthorKey      = KeyName("Author")
	VersionKey     = KeyName("Version")
)

func (n KeyName) Instances() []KeyName {
	return []KeyName{LanguagesKey, DescriptionKey, AuthorKey, VersionKey}
}

type KeyTarget string

const (
	GraphTarget = KeyTarget("graph")
	NodeTarget  = KeyTarget("node")
	EdgeTarget  = KeyTarget("edge")
)

func (t KeyTarget) Instances() []KeyTarget {
	return []KeyTarget{GraphTarget, NodeTarget, EdgeTarget}
}

type KeyType string

const (
	StringType = KeyType("string")
	IntType    = KeyType("int")
	XmlType    = KeyType("xml")
)

func (t KeyType) Instances() []KeyType {
	return []KeyType{StringType, IntType, XmlType}
}

type Key struct {
	ID     string
	Name   string
	Target KeyTarget
	Type   KeyType
}

func NewKey(key graphml.Key) (*Key, error) {
	target, err := enum.EnumByValue[KeyTarget](GraphTarget, KeyTarget(key.For))
	if err != nil {
		return nil, err
	}
	var t KeyType
	switch key.Type {
	case "string":
		t = StringType
	case "int":
		t = IntType
	case "":
		t = XmlType
	default:
		return nil, fmt.Errorf("type %s defined by key '%s' is not supported by this application", key.Type, key.Name)
	}
	return &Key{
		ID:     key.ID,
		Name:   key.Name,
		Target: *target,
		Type:   t,
	}, nil
}

type Keys []Key

func NewKeys[T KeyType](keys []graphml.Key) (*Keys, error) {
	var rsl Keys
	for _, key := range keys {
		if key.Name == "" {
			continue
		}
		r, err := NewKey(key)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *r)
	}
	return &rsl, nil
}

type Node struct {
	src graphml.Node
}

type Edge struct {
	src graphml.Edge
}

type Story struct {
	src    graphml.Graph
	Starts []Node
}

func NewStroy(doc graphml.Document) (*Story, error) {
	if len(doc.Graphs) == 0 {
		return nil, fmt.Errorf("story file doesn't contain a graph")
	} else if len(doc.Graphs) > 1 {
		return nil, fmt.Errorf("story file contains more than one graph (actual: %d)", len(doc.Graphs))
	}
	graph := doc.Graphs[0]
	fmt.Printf("%+v\n", doc.Keys)

	keys, err := NewKeys[KeyType](doc.Keys)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", keys)

	return &Story{
		src: graph,
	}, nil
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	doc, err := graphml.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	_, err = NewStroy(*doc)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%+v\n", story)
	// fmt.Printf("%s\n", doc.Graphs[0].Nodes[0].Data[0].Data[0])

	/*
		graph := graphml.NewGraphML("Story")
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		if err := graph.Decode(file); err != nil {
			log.Fatal(err)
		}
	*/

	// gui.Run()

}
