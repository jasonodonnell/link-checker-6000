package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"link-checker-6000/webcrawler"
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

	crawlerConfig := webcrawler.Config{
		PoolSize:       config.WorkerPool,
		MaxDepth:       config.MaxDepth,
		Timeout:        config.Timeout,
		AllowedDomains: config.AllowedDomains,
		DeniedDomains:  config.DeniedDomains,
		Logger:         logger,
	}

	crawler := webcrawler.NewWebCrawler(crawlerConfig)
	urls := []string{config.InitialURL}
	for len(urls) != 0 {
		docs := crawler.GetHTML(urls)

		urls = []string{}
		for _, doc := range docs {
			traverse(doc, config.BaseURL, &urls)
		}
	}

	fmt.Println("Dead links:")
	for _, dead := range crawler.DeadLinks() {
		fmt.Println(dead)
	}
}

func usage() {
	fmt.Println()
	fmt.Println("Usage: link-checker-6000 -config=/path/to/config.yaml")
}
