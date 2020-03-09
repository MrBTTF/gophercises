package sitemap

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
    "encoding/xml"

	"github.com/MrBTTF/gophercises/link"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}



func XML(url string) ([]byte, error) {
	root := &urlset{
		Xmlns: xmlns,
	}

	links := make(map[string]struct{})
	err := bfs(url, links)
	if err != nil {
		return nil, err
	}
	for link := range links{
		root.Urls = append(root.Urls, loc{link})
	}

	out, err := xml.MarshalIndent(root, "", "	")
	if err != nil {
		return nil, err
	}
	return out, nil
}


func getLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	// fmt.Println("Going to", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	links := link.ParseLinks(bytes.NewReader(body))

	var result []string
	for _, link := range links {
		if strings.Contains(link.URL, url) || link.URL[0] == '/' {
			result = append(result, link.URL)
		}
	}

	return result, err
}

func bfs(domain string, result map[string]struct{}) error {
	queue := []string{domain}
	visited := map[string]struct{}{}

	var f func() error
	f = func() error {
		if len(queue) == 0 {
			return nil
		}

		link := queue[0]
		if _, ok := visited[link]; ok {
			return nil
		}
		visited[link] = struct{}{}
		queue = queue[1:]

		children, err := getLinks(link)
		if err != nil {
			return err
		}
		// fmt.Println("Children", children)
		// fmt.Println("added")
		for _, child := range children {
			if _, ok := result[child]; !ok {
				if !strings.Contains(child, domain)  {
					child = domain + child
				}
				queue = append(queue, child)
				result[child] = struct{}{}
				// fmt.Printf("%s ", child)
			}
		}
		// fmt.Println("\nResult: ", result)
		// fmt.Println()

		return f()
	}

	return f()
}
