package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"path"
	"html/template"
)

type ArcHandler struct {
	arcs map[string]Arc
	mux http.Handler
}

func NewArcHandler(arcs map[string]Arc) *ArcHandler{
	return &ArcHandler{
		arcs: arcs,
	}
}

func (s *ArcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, arcName := path.Split(r.URL.Path)
	arc, ok := s.arcs[arcName]
	if !ok {
		http.Redirect(w, r, "/intro", 301)
	}
	t, err := template.ParseFiles("arc.html")
	if err != nil {
		fmt.Println(err)
	}
    t.Execute(w, &arc)
}


type Option struct {
	Text string
	NextArc string `json:"arc"`
}

type Arc struct {
	Title string
	Story []string
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

func main()  {
	fmt.Println("Choose Your Own Adventure server on port 8080")

	arcs, err := LoadBook("assets/gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("arcs loaded successfully")

	http.Handle("/",  NewArcHandler(arcs))
	http.Handle("/css/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":8080", nil)
}
