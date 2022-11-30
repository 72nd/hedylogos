package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/72nd/hedylogos/enum"
	"github.com/freddy33/graphml"
)

// The different key names as defined by the specification.
type KeyName string

const (
	LanguagesKey    = KeyName("Languages")
	DescriptionKey  = KeyName("Description")
	AuthorKey       = KeyName("Author")
	VersionKey      = KeyName("Version")
	TypeKey         = KeyName("Type")
	NodeGraphicsKey = KeyName("yd.nodegraphics")
)

func (n KeyName) Instances() []KeyName {
	return []KeyName{LanguagesKey, DescriptionKey, AuthorKey, VersionKey, TypeKey, NodeGraphicsKey}
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

type DataType = interface {
	string | int | xml.Token
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
	Type KeyType
}

// Returns a new Key from a graphml Key.
func NewKey(key graphml.Key) (*Key, bool, error) {
	target, err := enum.EnumByValue[KeyTarget](GraphTarget, KeyTarget(key.For))
	if err != nil {
		return nil, true, nil
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
		return nil, false, fmt.Errorf("type %s defined by key '%s' is not supported by this application", key.Type, key.Name)
	}
	if err != nil {
		return nil, false, err
	}
	name := keyName(key)
	if name == "" {
		return nil, true, nil
	}
	return &Key{
		ID:     key.ID,
		Name:   name,
		Target: *target,
		Type:   t,
	}, false, nil
}

// If the method returns an empty string t
func keyName(key graphml.Key) string {
	if key.Name != "" {
		return key.Name
	}
	if len(key.Object.Unrecognized) == 0 {
		return ""
	}
	if key.Object.Unrecognized[0].Name.Local == "yfiles.type" {
		return fmt.Sprintf("yd.%s", key.Object.Unrecognized[0].Value)
	}
	return ""
}

// Looks for a data entry matching the given key.
func (k Key) matchDataElement(obj graphml.ExtObject) (*graphml.Data, error) {
	for _, data := range obj.Data {
		if data.Key != k.ID {
			continue
		}
		return &data, nil
	}
	return nil, fmt.Errorf("object doesn't contain a data entry for %s", k.debugName())
}

// Looks for a matching data entry and extracts it's value.
func (k Key) matchAndGetValue(obj graphml.ExtObject) ([]DataType, error) {
	ele, err := k.matchDataElement(obj)
	if err != nil {
		return nil, err
	}
	if len(ele.Data) == 0 {
		return nil, fmt.Errorf("data entry matched by name %s contains no data", k.debugName())
	}
	if k.Type == XmlType {
		return k.parseDataAsXml(*ele)
	}
	charVal, ok := ele.Data[0].(xml.CharData)
	if ok {
		// Data is represented as `xml.CharData`.
		switch k.Type {
		case StringType:
			return []DataType{string(charVal)}, nil
		case IntType:
			v, err := strconv.Atoi(string(charVal))
			return []DataType{v}, err
		default:
			return nil, fmt.Errorf("key %s defines %s as type which isn't implemented by the Key class", k.debugName(), k.Type)
		}
	}
	_, ok = ele.Data[0].(xml.StartElement)
	if ok {
		switch k.Type {
		case StringType:
			return nil, fmt.Errorf("string data entry matched by %s is expected to be represented by chars in the source but XML data was found instead", k.debugName())
		case IntType:
			return nil, fmt.Errorf("int data entry matched by %s is expected to be represented by chars in the source but XML data was found instead", k.debugName())
		}
	}
	return nil, fmt.Errorf("logic error inside data entry matched by %s, this case isn't implemented", k.debugName())
}

// Parses a `graphml.Data` as an XML element.
func (k Key) parseDataAsXml(data graphml.Data) ([]DataType, error) {
	var rsl []DataType
	for _, ele := range data.Data {
		value, ok := ele.(xml.StartElement)
		if ok {
			rsl = append(rsl, value)
		}
		char, ok := ele.(xml.CharData)
		if ok && !strings.HasPrefix(string(char), "\n") {
			rsl = append(rsl, string(char))
		}
	}
	if len(rsl) == 0 {
		return nil, fmt.Errorf("xml data entry matched by %s doesn't contain any data", k.debugName())
	}
	return rsl, nil
}

// Returns a string containing the name and id of a key instance.
// Used in error messages.
func (k Key) debugName() string {
	return fmt.Sprintf("name '%s' (ID: %s)", k.Name, k.ID)
}

// A collection of multiple instances of Keys.
type Keys map[string]Key

// Creates a new Keys collection from a graphml Key slice. slice.
func NewKeys[T DataType](keys []graphml.Key) (*Keys, error) {
	var rsl = make(Keys)
	for _, key := range keys {
		r, ignore, err := NewKey(key)
		if ignore {
			continue
		}
		if err != nil {
			return nil, err
		}
		rsl[r.Name] = *r
	}
	return &rsl, nil
}

// Returns the data by a given key name.
func (k Keys) DataByName(obj graphml.ExtObject, name string) ([]DataType, error) {
	key, ok := k[name]
	if !ok {
		return nil, fmt.Errorf("there is no key with the name '%s'", name)
	}
	return key.matchAndGetValue(obj)
}
