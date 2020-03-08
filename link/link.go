package link

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"strings"
)

type Link struct {
	URL  string
	Text string
}

func nextNode(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		fmt.Println("Node:")
		fmt.Println(n.Data)
		fmt.Println(n.Type)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nextNode(c)
	}
}

func parseText(n *html.Node) string {
	var s string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s += parseText(c)
	}
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data) + s
	}
	return s
}

func ParseLinks(r io.Reader) []Link {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	var links []Link

	var collectLinks func(n *html.Node)
	collectLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var url string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url = attr.Val
					break
				}
			}
			if url == "" {
				return
			}
			links = append(links, Link{
				URL:  url,
				Text: parseText(n),
			})
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			collectLinks(c)
		}
	}
	collectLinks(doc)
	return links
}
