package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	WorkerPool     int      `koanf:"workerPool"`
	MaxDepth       int      `koanf:"maxDepth"`
	Timeout        int      `koanf:"timeout"`
	AllowedDomains []string `koanf:"allowedDomains"`
	DeniedDomains  []string `koanf:"deniedDomains"`
}

var k = koanf.New(".")

func loadConfig(path string) (*Config, error) {
	k.Load(file.Provider(path), yaml.Parser())
	out := &Config{}
	k.Unmarshal("", out)

	if out.WorkerPool == 0 {
		out.WorkerPool = 5
	}

	if out.MaxDepth == 0 {
		out.MaxDepth = 5
	}

	if out.Timeout == 0 {
		out.Timeout = 5
	}

	return out, nil
}

type opts struct {
	configPath string
	initialURL string
	baseURL    string
}

func parseFlags() opts {
	path := flag.String("config", "", "path to a config.yaml file")
	base := flag.String("baseURL", "", "the URL to use when rewriting relative links")
	flag.Usage = usage

	flag.Parse()
	if *path == "" {
		flag.Usage()
		os.Exit(2)
	}

	if *base == "" {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	return opts{
		configPath: *path,
		baseURL:    *base,
		initialURL: args[0],
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of link-checker-6000:\n")
	fmt.Fprintf(os.Stderr, "\tlink-checker-6000 [flags] [starting url]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}
