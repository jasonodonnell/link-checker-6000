package main

import (
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
