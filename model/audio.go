package model

import (
	"fmt"

	"github.com/72nd/hedylogos/graphml"
)

// Collection of multiple [Audio] instances.
type Audios []Audio

// Returns a new [Audios] instance based on a slice of graphml.Audio
// instances.
func NewAudios(audios []graphml.Audio, langs Languages) (*Audios, error) {
	var rsl Audios
	for _, audio := range audios {
		a, err := NewAudio(audio, langs)
		if err != nil {
			return nil, err
		}
		rsl = append(rsl, *a)
	}
	return &rsl, nil
}

// Custom type. An output type using audio files.
type Audio struct {
	// Path to audio file.
	Source string
	// The ID [Language.ID] of the language which the audio file is in.
	Language string
	// If the content of the audio file mentions a number key this can
	// be noted in this field. This information is only used for
	// debugging and has therefore no direct technical function.
	Target string
	// Textual representation of the content of the audio file. This
	// is used to generate a script for producing the audio files in
	// the first place.
	Text string
}

// Returns a new Audio instance based on a graphml.Audio instance.
func NewAudio(audio graphml.Audio, langs Languages) (*Audio, error) {
	if audio.Language != "" && audio.Language != "*" && !langs.Exists(audio.Language) {
		return nil, fmt.Errorf("audio file with pad '%s' defines a language (%s) which isn't declared within the graph", audio.Source, audio.Language)
	}
	return &Audio{
		Source:   audio.Source,
		Language: audio.Language,
		Target:   audio.Target,
		Text:     audio.Text,
	}, nil
}

func (a Audio) Execute() {
	fmt.Printf("Play file '%s'", a.Source)
}
