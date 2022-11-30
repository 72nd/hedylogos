package model

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/72nd/hedylogos/graphml"
)

// Collection of [Language] instances.
type Languages []Language

// Takes a slice of [github.com/72nd/hedylogos/graphml.Language] and returns
// a new [Languages] instance.
func NewLanguages(langs []graphml.Language) (*Languages, error) {
	var rsl Languages
	for _, lang := range langs {
		l, err := NewLanguage(lang)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *l)
	}
	return &rsl, nil
}

// Custom type. Defines the available languages in the scenario. Each
// node with a language specific functions (like audio) has to
// provide an implementation for all languages.
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

// Takes and validates an instance of [github.com/72nd/hedylogos/graphml.Language]
// and returns a new [Language] instance.
func NewLanguage(lang graphml.Language) (*Language, error) {
	idExp, err := regexp.Compile(`^(\w{2}|\w{2}_\w{2})$`)
	if err != nil {
		return nil, err
	}
	if !idExp.MatchString(lang.ID) {
		return nil, fmt.Errorf("given language id '%s' has an invalid format (ex. for allowed de or de_ch)", lang.ID)
	}
	if lang.Name == "" {
		return nil, fmt.Errorf("name of given language (with id '%s') is empty", lang.ID)
	}
	value, err := strconv.Atoi(lang.Value)
	if err != nil {
		return nil, fmt.Errorf("given language (with id '%s') value '%s' cannot be converted to an integer", lang.ID, lang.Value)
	}
	if value < 0 || value > 9 {
		return nil, fmt.Errorf("language (with '%s') value has to between 0 and 9 but is %d", lang.ID, value)
	}
	return &Language{
		ID:    lang.ID,
		Name:  lang.Name,
		Value: value,
	}, nil
}
