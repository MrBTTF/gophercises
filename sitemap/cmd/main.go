package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/MrBTTF/gophercises/sitemap"
)

func main() {
	urlFlag := flag.String("url", "go get ", "Website URL to build a sitemap of")
	flag.Parse()

	xml, err := sitemap.XML(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("sitemap.xml", xml, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
