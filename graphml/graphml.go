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
	Default *Default `xml:"default"`
}

type Default struct {
	XMLName  xml.Name `xml:"default"`
	XMLSpace string   `xml:"xml:space,attr"`
	Body     string   `xml:",chardata"`
}

// Represents a GraphML Document.
type Document struct {
	XMLName xml.Name `xml:"graphml"`
	Keys    []Key    `xml:"key"`
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
