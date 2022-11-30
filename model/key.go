package model

import (
	"fmt"

	"github.com/72nd/hedylogos/enum"
	"github.com/72nd/hedylogos/graphml"
)

// A collection of multiple instances of Keys.
type Keys []Key

// Creates a new Keys collection from a graphml Key slice. slice.
func NewKeys(keys []graphml.Key) (*Keys, error) {
	var rsl Keys
	for _, key := range keys {
		k, err := NewKey(key)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *k)
	}
	return &rsl, nil
}

// Returns a key with the given id. Returns an error if no key was found.
func (k Keys) ByID(id string) (*Key, error) {
	for _, key := range k {
		if key.ID == id {
			return &key, nil
		}
	}
	return nil, fmt.Errorf("no key was found for id '%s'", id)
}

// Returns all duplicate names of the keys.
func (k Keys) Duplicates() []string {
	var rsl []string
	for i := range k {
		for j := range k {
			if i == j {
				continue
			}
			if k[i].Name == k[j].Name {
				rsl = append(rsl, k[i].Name)
			}
		}
	}
	return rsl
}

// Defines the targets for any key.
type KeyTarget string

const (
	GraphTarget   = KeyTarget("graph")
	NodeTarget    = KeyTarget("node")
	EdgeTarget    = KeyTarget("edge")
	PortTarget    = KeyTarget("port")
	GraphMLTarget = KeyTarget("graphml")
)

func (t KeyTarget) Instances() []KeyTarget {
	return []KeyTarget{GraphTarget, NodeTarget, EdgeTarget, PortTarget, GraphMLTarget}
}

// Data type used by the source value of the key.
type SourceKeyType string

const (
	StringSourceType = SourceKeyType("string")
	IntSourceType    = SourceKeyType("int")
	XmlSourceType    = SourceKeyType("xml")
)

func (t SourceKeyType) Instances() []SourceKeyType {
	return []SourceKeyType{StringSourceType, IntSourceType, XmlSourceType}
}

// A key is used to access a certain data attribute of a Graph, Node or Edge.
type Key struct {
	// Unique ID designated by the graphml file. Each data field uses this
	// ID to refer to a certain key. The ID can differ between different
	// versions of a graphml file.
	ID string
	// User given name of the key. While the ID is used to link data fields
	// with keys within a certain graphml file the name always describes the
	// same data field type.
	Name string
	// States whether the key applies to a data attribute of a Graph, Node or
	// Edge.
	Target KeyTarget
	// States the datatype of a data field linked to the key.
	Type SourceKeyType
}

// Returns a new Key from a graphml Key. Second return value states whether
// the key should be ignored or not.
func NewKey(key graphml.Key) (*Key, error) {
	target, err := enum.EnumByValue[KeyTarget](GraphTarget, KeyTarget(key.For))
	if err != nil {
		return nil, fmt.Errorf("error while importing key '%s' (ID: %s), %s", key.Name, key.ID, err)
	}
	tp, err := enum.EnumByValue[SourceKeyType](StringSourceType, SourceKeyType(key.Type))
	if err != nil {
		// If no type is given by the source file it's assumed to be [XmlType].
		if key.Type != "" {
			return nil, err
		}
		tmp := XmlSourceType
		tp = &tmp
	}
	name := key.Name
	if key.YFileType != "" {
		name = fmt.Sprintf("yft.%s", key.YFileType)
	}
	return &Key{
		ID:     key.ID,
		Name:   name,
		Target: *target,
		Type:   *tp,
	}, nil
}
