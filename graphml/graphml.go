// This package provides a very simple importer for a GraphML file.
// It tries to unmarshal all data field values custom to hedylogos
// as defined in the [input file specification] it also takes some
// fields specific to yEd into account.
//
// Thus this module shouldn't be used on other GraphML files as the
// sole propose of this package is to provide the raw data for
// further processing within hedylogos as structured as possible.
// Because of this reason there are no type-checks whatsoever.
//
// You can learn more about the GraphML format by reading it's
// [specification]. Most documentation of structs within this package
// is also taken from this source.
//
// [input file specification]: https://github.com/72nd/hedylogos/blob/master/doc/specification.md
// [specification]: http://graphml.graphdrawing.org/specification.html
package graphml

import (
	"encoding/xml"
	"os"
)

// GraphML type. GraphML provide a mechanism to add data to the
// structural elements (e.g. `<graph>`s, `<node>`s, `<edge>`s, etc.).
// Data labellings are considered to be (possibly partial functions
// that assign values in an (a priori) arbitrary range to elements
// of the graph. Such a function is declared by a `<key>` element.
// The domain of definition of this function is specified by the
// 'for' attribute of the `<key>“.
type Key struct {
	XMLName xml.Name `xml:"key"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"attr.name,attr"`
	Type    string   `xml:"attr.type,attr"`
	For     string   `xml:"for,attr"`
	Default string   `xml:"default"`
}

// Custom type. Defines the available languages in the scenario. Each
// node with a language specific functions (like audio) has to
// provide an implementation for all languages.
type Language struct {
	XMLName xml.Name `xml:"language"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

// yEd type. Contains information about the visual representation of
// a given node. As some information (label text and shape) are used
// within hedylogos the importer provides a partially implementation
// of this type.
type ShapeData struct {
	XMLName xml.Name `xml:"ShapeNode"`
	Label   string   `xml:"NodeLabel"`
	Shape   Shape    `xml:"Shape"`
}

// yEd type. Contains information about the geometric shape of the
// node. This information is used withing hedylogos to check if the
// input files follows the specification in regard to the correct
// usage of the available shapes in relation to the node type.
//
// The only reason this is represented as it's own separate struct
// is that it's not possible to use the `,attr` notation of an xml
// struct-tag in conjunction with an `>` expression.
type Shape struct {
	XMLName xml.Name `xml:"Shape"`
	Type    string   `xml:"type,attr"`
}

// GraphML type. GraphML provide a mechanism to add data to the
// structural elements (e.g. `<graph>`s, `<node>`s, `<edge>`s, etc.).
// Data labellings are considered to be (possibly partial) functions
// that assign values in an (a priori) arbitrary range to elements of
// the graph. Values are defined by a `<default>` element (child of `<key>`)
// and/or `<data>` elements (children of the elements, which are in the
// domain of definition) whose 'key' attribute-values matches the 'id'
// of the `<key>`
type Data struct {
	XMLName   xml.Name   `xml:"data"`
	Key       string     `xml:"key,attr"`
	XMLSpace  string     `xml:"xml:space,attr"`
	Languages []Language `xml:"language"`
	ShapeData ShapeData  `xml:"ShapeNode"`
	Value     string     `xml:",innerxml"`
}

// GraphML type. A node within a graph.
type Node struct {
	XMLName xml.Name `xml:"node"`
	ID      string   `xml:"id,attr"`
	Data    []Data   `xml:"data"`
}

// GraphML type. A edge within a graph connecting two nodes.
type Edge struct {
	XMLName xml.Name `xml:"edge"`
	ID      string   `xml:"id,attr"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
	Data    []Data   `xml:"data"`
}

// GraphML type. Defines a Graph within a GraphML document.
type Graph struct {
	XMLName     xml.Name `xml:"graph"`
	ID          string   `xml:"id,attr"`
	EdgeDefault string   `xml:"edgedefault,attr"`
	Data        []Data   `xml:"data"`
	Nodes       []Node   `xml:"node"`
	Edges       []Edge   `xml:"edge"`
}

// Represents a GraphML Document.
type Document struct {
	XMLName xml.Name `xml:"graphml"`
	Keys    []Key    `xml:"key"`
	Graph   Graph    `xml:"graph"`
}

// Loads a GraphML file with the given path and tries to parse it
// into a `Document`.
func FromFile(path string) (*Document, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rsl Document
	xml.Unmarshal(file, &rsl)
	return &rsl, nil
}
