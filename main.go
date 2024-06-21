package main

import (
	"fmt"
	"strings"
	"time"

	"link-checker-6000/fetcher"

	"golang.org/x/net/html"
)

func extractURLs(n *html.Node, base string) []string {
	var urls []string
	f := func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "html") {
			for _, a := range n.Attr {
				if a.Key == "href" {
					value := a.Val
					if a.Val[0] == '/' {
						value = base + a.Val
					}

					if strings.Contains(a.Val, base) {
						urls = append(urls, value)
					}
					break
				} else if a.Key == "path" {
					value := a.Val
					if a.Val[0] == '/' {
						value = base + a.Val
					}

					if strings.Contains(a.Val, base) {
						urls = append(urls, value)
					}
					break
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		f(c)
	}
	f(n)
	return urls
}

func main() {
	urls := []string{"https://developer.hashicorp.com/vault/docs"}
	baseURL := "https://www.developer.hashicorp.com/vault/docs"

	f := fetcher.NewFetcher(urls, 5, 10, time.Second*5)
	for urls != nil {
		docs := f.GetHTML(urls)

		var results []string
		for _, doc := range docs {
			results = extractURLs(&doc, baseURL)
		}
		urls = results
	}

	fmt.Println(f.DeadLinks())
}
