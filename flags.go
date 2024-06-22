package main

import (
	"flag"
	"fmt"
	"os"
)

type opts struct {
	configPath string
	initialURL string
	baseURL    string
}

func parseFlags() opts {
	path := flag.String("config", "", "path to a config.yaml file")
	flag.Usage = usage

	flag.Parse()
	if *path == "" {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) <= 1 {
		flag.Usage()
		os.Exit(2)
	}

	return opts{
		configPath: *path,
		initialURL: args[0],
		baseURL:    args[1],
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of link-checker-6000:\n")
	fmt.Fprintf(os.Stderr, "\tlink-checker-6000 [flags] [url]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}
