package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ohmydevops/arvancloud-radar-notif/radar"
)

type Config struct {
	Service      string
	ShowServices bool
}

func parseFlags() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Service, "service", "", "Service name to monitor (e.g. google, github, etc.)")
	flag.BoolVar(&cfg.ShowServices, "services", false, "Show list of available services")

	// Save default flag.Usage so we can call it
	defaultUsage := flag.Usage
	flag.Usage = func() {
		defaultUsage() // print default usage
		fmt.Println("\nAvailable services:")
		for _, s := range radar.AllServices {
			fmt.Printf("  - %s\n", s)
		}
	}

	flag.Parse()

	// normalize to lowercase
	cfg.Service = strings.ToLower(cfg.Service)

	// Validate service
	if _, ok := radar.ParseService(cfg.Service); !ok {
		return nil, fmt.Errorf("invalid service: %s", cfg.Service)
	}
	return &cfg, nil
}
