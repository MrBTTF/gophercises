package sitemap

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MrBTTF/gophercises/link"
)

func XML(url string) ([]byte, error) {
	links, err := getLinks(url)
	if err != nil {
		return nil, err
	}

	fmt.Println(links)
	return nil, nil
}

func getLinks(url string) ([]link.Link, error) {
	resp, err := http.Get(url)
	fmt.Println(url)
	url = strings.ReplaceAll(url, "http://", "")
	url = strings.ReplaceAll(url, "https://", "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	links := link.ParseLinks(bytes.NewReader(body))

	var result []link.Link
	for _, link := range links {
		fmt.Println(link.URL)
		if strings.Contains(link.URL, url) || link.URL[0] == '/' {
			result = append(result, link)
		}
	}

	return result, err
}
