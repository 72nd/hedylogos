package oldmodel

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/freddy33/graphml"
)

// Represents a language which a story supports.
type Language struct {
	// Unique two-letter identifier for the language according to ISO 639-1.
	// If you're using a local variation of an language (dialect) use the
	// language code followed by an `_` and a two letter code of the country
	// region the local variation originates for (example `de_ch` for swiss
	// german).
	ID string
	// Name of the language.
	Name string
	// A number between 0 and 9 which is used to choose the language in an
	// language menu. Make sure to assign a unique number to each language
	// in a story.
	Value int
}

// Returns a new instance of `Language` based on a given XML element.
func NewLanguage(ele xml.Token) (*Language, error) {
	val, ok := ele.(xml.StartElement)
	if !ok {
		return nil, fmt.Errorf("a xml.StartElement is expected when parsing a language element got %T instead", ele)
	}
	if val.Name.Local != "language" {
		return nil, fmt.Errorf("a language has to be defined by an XML element with the name 'language' got '%s' instead", val.Name.Local)
	}
	var id string
	var name string
	value := -1
	for _, atr := range val.Attr {
		switch atr.Name.Local {
		case "id":
			id = atr.Value
		case "name":
			name = atr.Value
		case "value":
			v, err := strconv.Atoi(atr.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse value tag of language element %s", err)
			}
			value = v
		}
	}
	if id == "" {
		return nil, fmt.Errorf("failed to parse language element, id is not set")
	}
	if name == "" {
		return nil, fmt.Errorf("failed to parse language element, name is not set")
	}
	if value == -1 {
		return nil, fmt.Errorf("failed to parse language element, value is not set")
	}
	return &Language{
		ID:    id,
		Name:  name,
		Value: value,
	}, nil
}

// Contains all languages supported by the story.
type Languages []Language

// Returns a new instance of Languages according to the ExtObject of the
// graph.
func NewLanguages(obj graphml.ExtObject, keys Keys) (Languages, error) {
	data, err := keys.DataByName(obj, string(LanguagesKey))
	if err != nil {
		return nil, err
	}
	var rsl Languages
	for _, ele := range data {
		lang, err := NewLanguage(ele)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *lang)
	}
	return rsl, nil
}

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
}

// Takes a GraphML Document and parses it into a Story instance. The input
// file has to follow the specification as given under `doc/specification.md`.
func NewStroy(doc graphml.Document) (*Story, error) {
	if len(doc.Graphs) == 0 {
		return nil, fmt.Errorf("story file doesn't contain a graph")
	} else if len(doc.Graphs) > 1 {
		return nil, fmt.Errorf("story file contains more than one graph (actual: %d)", len(doc.Graphs))
	}
	graph := doc.Graphs[0]
	keys, err := NewKeys[DataType](doc.Keys)
	if err != nil {
		return nil, err
	}
	version, err := keys.DataByName(graph.ExtObject, string(VersionKey))
	if err != nil {
		return nil, err
	}
	tVersion, ok := version[0].(string)
	if !ok {
		return nil, fmt.Errorf("version data has to be a string")
	}
	author, err := keys.DataByName(graph.ExtObject, string(AuthorKey))
	if err != nil {
		return nil, err
	}
	tAuthor, ok := author[0].(string)
	if !ok {
		return nil, fmt.Errorf("author data has to be a string")
	}
	desc, err := keys.DataByName(graph.ExtObject, string(DescriptionKey))
	if err != nil {
		return nil, err
	}
	tDesc, ok := desc[0].(string)
	if !ok {
		return nil, fmt.Errorf("description data has to be a string")
	}
	languages, err := NewLanguages(graph.ExtObject, *keys)
	if err != nil {
		return nil, err
	}

	nodes, err := NewNodes(graph, *keys)
	if err != nil {
		return nil, err
	}

	return &Story{
		Version:     tVersion,
		Author:      tAuthor,
		Description: tDesc,
		Languages:   languages,
		Nodes:       *nodes,
	}, nil
}
