package main

import (
	"fmt"
	"log/slog"
	"os"

	"link-checker-6000/webcrawler"

	"golang.org/x/net/html"
)

func main() {
	logger := slog.Default()
	opts := parseFlags()

	config, err := loadConfig(opts.configPath)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(2)
	}

	crawlerConfig := webcrawler.Config{
		PoolSize:       config.WorkerPool,
		MaxDepth:       config.MaxDepth,
		Timeout:        config.Timeout,
		AllowedDomains: config.AllowedDomains,
		DeniedDomains:  config.DeniedDomains,
		Logger:         logger,
	}

	crawler := webcrawler.NewWebCrawler(crawlerConfig)
	urls := []string{opts.initialURL}

	for len(urls) != 0 {
		docs := crawler.GetHTML(urls)

		urls = []string{}
		for _, doc := range docs {
			traverse(doc, opts.baseURL, &urls)
		}
	}

	fmt.Println("Dead links:")
	for _, dead := range crawler.DeadLinks() {
		fmt.Println(dead)
	}
}

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
