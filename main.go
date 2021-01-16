package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type crawler struct {
	wg      sync.WaitGroup
	m       sync.Mutex
	visited map[string]struct{}
}

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

	c := crawler{
		visited: make(map[string]struct{}),
	}

	<-c.Start(context.Background(), f.rawURL)
}

func (c *crawler) Start(ctx context.Context, url string) chan struct{} {
	done := make(chan struct{})
	defer close(done)
	c.wg.Add(1)
	go c.fetch(ctx, url, url)
	c.wg.Wait()
	return done
}

func (c *crawler) fetch(ctx context.Context, url, baseURL string) {
	defer c.wg.Done()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		unhandled(err)
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		unhandled(err)
	}

	_ = parse(res)
}

func parse(res *http.Response) []string {
	fmt.Println(res.Header.Get("Content-Type"))
	return nil
}

func (f flags) validate() error {
	if len(f.path) == 0 {
		return fmt.Errorf("path is required")
	}
	fi, err := os.Lstat(f.path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("path is not a direcotry")
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

func unhandled(err error) {
	panic(fmt.Errorf("unhandled error: %w", err))
}
