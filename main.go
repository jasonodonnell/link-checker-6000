package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"link-checker-6000/fetcher"
)

func main() {
	logger := slog.Default()

	path := flag.String("config", "", "The path to a config.yaml file")
	flag.Parse()
	if *path == "" {
		logger.Error("config flag not set")
		usage()
		os.Exit(2)
	}

	config, err := loadConfig(*path)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(2)
	}

	fetcher := fetcher.NewFetcher(config.WorkerPool, config.MaxDepth, config.Timeout, config.AllowedDomains, config.DeniedDomains, logger)
	urls := []string{config.InitialURL}
	for len(urls) != 0 {
		docs := fetcher.GetHTML(urls)

		urls = []string{}
		for _, doc := range docs {
			traverse(doc, config.BaseURL, &urls)
		}
	}

	fmt.Println("Dead links:")
	for _, dead := range fetcher.DeadLinks() {
		fmt.Println(dead)
	}
}

func usage() {
	fmt.Println()
	fmt.Println("Usage: link-checker-6000 -config=/path/to/config.yaml")
}
