package main

import (
	"fmt"
	"log"

	"github.com/MrBTTF/gophercises/cyoa"
)

func main() {
	fmt.Println("[INFO] Welcome to \"Choose Your Own Adventure\"!")

	arcs, err := cyoa.LoadBook("assets/gopher.json")
	if err != nil {
		log.Fatal(err)
	}


	arc := arcs["intro"]
	for {
		if len(arc.Options) == 0 {
			break
		}

		printArc(arc)

		var option int
		_, err = fmt.Scanf("%d\n", &option)

		if err != nil {
			var input string
			_, err := fmt.Scanf("%s\n", &input)
			if err != nil {
				log.Fatal(err)
			}
			if input == "q" {
				break
			}
			fmt.Printf("Unknown input %s\n", input)
		}
		if option >= len(arc.Options) {
			fmt.Printf("Please choose options in range 0-%d\n", len(arc.Options)-1)
		}
		arc = arcs[arc.Options[option].NextArc]
	}

	fmt.Printf("------%s------\n", arc.Title)
	for _, line := range arc.Story {
		fmt.Println("\t", line)
	}
	fmt.Println("-------THE END------")
}


func printArc(arc cyoa.Arc) {
	fmt.Printf("------%s------\n", arc.Title)
	for _, line := range arc.Story {
		fmt.Println("\t", line)
	}
	fmt.Println()
	for i, option := range arc.Options {
		fmt.Printf("%d: %s\n", i, option.Text)
	}
	fmt.Println()
	fmt.Println("q: Quit")
}
