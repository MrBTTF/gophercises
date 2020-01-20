package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type Option struct {
	Text    string
	NextArc string `json:"arc"`
}

type Arc struct {
	Title   string
	Story   []string
	Options []Option
}

func LoadBook(filename string) (map[string]Arc, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var arcs map[string]Arc
	json.Unmarshal(jsonData, &arcs)

	return arcs, nil
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	fmt.Printf("[INFO] \"Choose Your Own Adventure\" server is running on port %s\n", port)

	arcs, err := LoadBook("assets/gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[INFO] Story's arcs loaded successfully")

	http.Handle("/", arcHandler(arcs))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":"+port, nil)
}

func arcHandler(arcs map[string]Arc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, arcName := path.Split(r.URL.Path)
		arc, ok := arcs[arcName]
		if !ok {
			http.Redirect(w, r, "/intro", 301)
		}
		t, err := template.ParseFiles("assets/arc.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, &arc)
	})
}
