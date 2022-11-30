package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/72nd/hedylogos/graphml"
)

type DataValueType interface {
	string | int | Languages | graphml.ShapeData | Audios
}

func ValueByName[T DataValueType](sto Storage, name string) (*T, error) {
	entry, err := sto.ByKeyName(name)
	if err != nil {
		return nil, err
	}
	rsl, ok := entry.Value.(T)
	if ok {
		return &rsl, nil
	}
	ptRsl, ok := entry.Value.(*T)
	if ok {
		return ptRsl, nil
	}
	return nil, fmt.Errorf("error while trying to cast data entry %s (id: %s) to type %T", entry.Key.Name, entry.Key.ID, new(T))
}

// A collection of Data instances. Used only during the construction
// of the models from a GraphML file.
type Storage []Data

// A new [Storage] from a slice of input data instances.
func NewStorage(data []graphml.Data, keys Keys, langs *Languages) (*Storage, error) {
	var rsl Storage
	for _, dt := range data {
		d, err := NewData(dt, keys, langs)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *d)
	}
	return &rsl, nil
}

// Returns the data instance with a given key name. Returns an error if
// no matching [Data] entry was found or there is more than one entry
// for the same name.
func (s Storage) ByKeyName(name string) (*Data, error) {
	var rsl Data
	found := false
	for _, data := range s {
		if data.Key.Name == name {
			if found {
				return nil, fmt.Errorf("there is more than one Data entry with the same key name '%s'", name)
			}
			found = true
			rsl = data
		}
	}
	if !found {
		return nil, fmt.Errorf("no Data entry found for key name '%s'", name)
	}
	return &rsl, nil
}

// Represents a data entry within the document. Used only during the
// construction of the model while opening a file.
type Data struct {
	Key   Key
	Value any
}

// Parses a given [github.com/72nd/hedylogos/graphml.Data] instance.
func NewData(data graphml.Data, keys Keys, langs *Languages) (*Data, error) {
	key, err := keys.ByID(data.Key)
	if err != nil {
		return nil, err
	}
	if len(data.Languages) != 0 && key.Type == XmlSourceType {
		langs, err := NewLanguages(data.Languages)
		if err != nil {
			return nil, err
		}
		return &Data{
			Key:   *key,
			Value: langs,
		}, nil
	}
	if len(data.Audio) != 0 && key.Type == XmlSourceType {
		if langs == nil {
			return nil, fmt.Errorf("need languages to construct audio data fields")
		}
		audios, err := NewAudios(data.Audio, *langs)
		if err != nil {
			return nil, err
		}
		return &Data{
			Key:   *key,
			Value: audios,
		}, nil
	}
	if data.ShapeData.Label != "" && key.Type == XmlSourceType {
		return &Data{
			Key:   *key,
			Value: data.ShapeData,
		}, nil
	}
	strVal := strings.TrimSpace(data.CharData)
	if key.Type == IntSourceType {
		value, err := strconv.Atoi(strVal)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse content ('%s') of data '%s' (id: %s) as int value", strVal, key.Name, key.ID)
		}
		return &Data{
			Key:   *key,
			Value: value,
		}, nil
	}
	if key.Type == StringSourceType {
		return &Data{
			Key:   *key,
			Value: strVal,
		}, nil
	}
	return &Data{}, nil
}
