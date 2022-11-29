package main

import (
	"fmt"
	"log"
	"os"

	"github.com/72nd/nohelpline/model"
	"github.com/freddy33/graphml"
)

func main() {
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
	// fmt.Printf("%s\n", doc.Graphs[0].Nodes[0].Data[0].Data[0])

	/*
		graph := graphml.NewGraphML("Story")
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		if err := graph.Decode(file); err != nil {
			log.Fatal(err)
		}
	*/

	// gui.Run()

}
