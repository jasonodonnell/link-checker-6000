package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"link-checker-6000/fetcher"
)

func main() {
	logger := slog.Default()

	config, err := LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(2)
	}

	urls := []string{config.InitialURL}

	fetcher := fetcher.NewFetcher(urls, config.WorkerPool, config.MaxDepth, config.AllowedDomains, config.DeniedDomains)
	for len(urls) != 0 {
		docs := fetcher.GetHTML(context.Background, urls)
		urls = []string{}

		for _, doc := range docs {
			traverse(doc, config.BaseURL, &urls)
		}
	}

	fmt.Println(fetcher.DeadLinks())
}
