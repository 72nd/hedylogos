package main

import (
	"fmt"
	"log"
	"os"

	"github.com/72nd/hedylogos/graphml"
	"github.com/72nd/hedylogos/model"
)

func main() {
	/*
		// USING THE 3RD-PARTY LIBRARY
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		doc, err := graphml.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		story, err := model.NewStroy(*doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", story)
	*/

	// USING OWN GRAPHML LIBRARY
	doc, err := graphml.FromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	story, err := model.NewStory(*doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", story)
}
