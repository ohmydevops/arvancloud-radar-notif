package flag

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/ohmydevops/arvancloud-radar-notif/radar"
)

const DefaultCheckDelay = 1

type Config struct {
	Service      string
	ShowServices bool
	CheckDelay   int // In minutes
}

func ParseFlags() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Service, "service", "", "Service name to monitor (e.g. google, github, etc.)")
	flag.BoolVar(&cfg.ShowServices, "services", false, "Show list of available services")
	flag.IntVar(&cfg.CheckDelay, "delay", DefaultCheckDelay, "Delay between checks in minutes")

	flag.Parse()

	if cfg.ShowServices {
		return &cfg, nil
	}

	// normalize to lowercase
	cfg.Service = strings.ToLower(cfg.Service)

	if cfg.Service == "" {
		return nil, errors.New("must specify a service")
	}
	// Validate service
	if _, ok := radar.ParseService(cfg.Service); !ok {
		return nil, fmt.Errorf("invalid service: %s", cfg.Service)
	}

	if cfg.CheckDelay < 1 {
		return nil, fmt.Errorf("delay must be greater than 0")
	}
	return &cfg, nil
}
