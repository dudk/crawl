package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

type flags struct {
	rawURL string
	path   string
}

func main() {
	var f flags
	flag.StringVar(&f.rawURL, "url", "", "base URL for crawling")
	flag.StringVar(&f.path, "path", "", "path to save crawling result")
	flag.Parse()
	if err := f.validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid input: %s\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func (f flags) validate() error {
	if len(f.path) == 0 {
		return fmt.Errorf("path is required")
	}

	if len(f.rawURL) == 0 {
		return fmt.Errorf("URL is required")
	}

	u, err := url.Parse(f.rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if len(u.Scheme) == 0 {
		return fmt.Errorf("URL scheme is required")
	}
	if len(u.Host) == 0 {
		return fmt.Errorf("URL host is required")
	}

	return nil
}
