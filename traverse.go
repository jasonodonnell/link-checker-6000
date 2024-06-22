package main

import "golang.org/x/net/html"

func traverse(n *html.Node, base string, results *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				value := attr.Val
				if attr.Val[0] == '/' {
					value = base + value
				}
				*results = append(*results, value)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, base, results)
	}
}
