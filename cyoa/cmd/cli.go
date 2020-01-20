package main

import (
	"fmt"
	"log"

	"github.com/MrBTTF/gophercises/cyoa/cyoa"
)

func main() {
	fmt.Println("[INFO] Welcome to \"Choose Your Own Adventure\"!")

	arcs, err := cyoa.LoadBook("assets/gopher.json")
	if err != nil {
		log.Fatal(err)
	}

}
