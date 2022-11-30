package main

import (
	"fmt"
	"log"
	"os"

	"github.com/72nd/nohelpline/graphml"
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
	fmt.Println(doc)

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
