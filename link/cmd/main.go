package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MrBTTF/gophercises/link"
)

func readDir(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string

	fmt.Println("Found files")
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".html") {
			continue
		}
		fmt.Println(f.Name())

		data, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			return nil, err
		}

		result = append(result, string(data))
	}
	return result, nil
}

func parseFile(filename string) error{
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	links := link.ParseLinks(file)
	for _, link := range links {
		fmt.Println("URL: ", link.URL)
		fmt.Println("Text: ", link.Text)
	}
	return nil
}

func parseDir(dir string) error{
	files, err := readDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		links := link.ParseLinks(strings.NewReader(f))
		for _, link := range links {
			fmt.Println("URL: ", link.URL)
			fmt.Println("Text: ", link.Text)
		}
	}
	return nil
}

func main() {
	var file = flag.String("f", "", "html file")
	var dir = flag.String("d", "examples", "dir with html files")
	flag.Parse()

	fmt.Println(*file)
	if *file != "" {
		err := parseFile(*file)
		if err != nil {
			log.Fatal(err)
		}
		return
	} 

	err := parseDir(*dir)
	if err != nil {
		log.Fatal(err)
	}

}
