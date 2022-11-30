// This package provides a very simple importer for any gr
package graphml

import (
	"encoding/xml"
	"os"
)

type Key struct {
	XMLName xml.Name `xml:"key"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"attr.name,attr"`
	Type    string   `xml:"attr.type,attr"`
	For     string   `xml:"for,attr"`
	Default string   `xml:"default"`
}

type Language struct {
	XMLName xml.Name `xml:"language"`
	Name    string   `xml:"name,attr"`
}

type ShapeNode struct {
	XMLName xml.Name `xml:"ShapeNode"`
	Label   string   `xml:"NodeLabel"`
}

type ShapeData struct {
	XMLName xml.Name `xml:"ShapeNode"`
}

type Data struct {
	XMLName   xml.Name   `xml:"data"`
	Key       string     `xml:"key,attr"`
	XMLSpace  string     `xml:"xml:space,attr"`
	Languages []Language `xml:"language"`
	ShapeData ShapeData  `xml:"ShapeNode"`
	Value     string     `xml:",innerxml"`
}

type Graph struct {
	XMLName     xml.Name `xml:"graph"`
	ID          string   `xml:"id,attr"`
	EdgeDefault string   `xml:"edgedefault,attr"`
	Data        []Data   `xml:"data"`
	Nodes       []Node   `xml:"node"`
	Edges       []Edge   `xml:"edge"`
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	ID      string   `xml:"id,attr"`
	Data    []Data   `xml:"data"`
}

type Edge struct {
	XMLName xml.Name `xml:"edge"`
	ID      string   `xml:"id,attr"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
	Data    []Data   `xml:"data"`
}

// Represents a GraphML Document.
type Document struct {
	XMLName xml.Name `xml:"graphml"`
	Keys    []Key    `xml:"key"`
	Graph   Graph    `xml:"graph"`
}

func FromFile(path string) (*Document, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rsl Document
	xml.Unmarshal(file, &rsl)
	return &rsl, nil
}
