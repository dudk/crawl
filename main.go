package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/dudk/crawl/caching"
	"github.com/dudk/crawl/fetch"
	"github.com/dudk/crawl/parsing"
	"github.com/dudk/crawl/scheduling"
)

type flags struct {
	rawURL string
}

func main() {
	var f flags
	flag.StringVar(&f.rawURL, "url", "", "Fully-specified base URL for crawling. If it's not ended with a slash, single URL will be fetched.")
	flag.Parse()
	if err := f.validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid input: %s\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	s := scheduling.UnboundScheduler{
		BaseURL: f.rawURL,
		ErrorHandler: scheduling.ErrorHandlerFunc(func(err error) {
			fmt.Fprintf(os.Stdout, "fetcher error: %s\n", err.Error())
		}),
		Fetcher: fetch.Fetcher{
			Client: &http.Client{
				Timeout: 10 * time.Second,
			},
			Parser: parsing.HTML{
				BaseURL: f.rawURL,
				Visitor: caching.NewInMemoryCache(),
			},
			BodyReader: fetch.BodyReaderFunc(func(ctx context.Context, URL *url.URL, r io.Reader) error {
				fmt.Printf("fetched: %v\n", URL)
				if _, err := io.Copy(ioutil.Discard, r); err != nil {
					return fmt.Errorf("error discarding body: %w", err)
				}
				return nil
			}),
		},
	}
	ctx, cancelFn := context.WithCancel(context.Background())

	sigint := make(chan os.Signal, 1)
	// interrupt signal
	signal.Notify(sigint, os.Interrupt)
	go func() {
		// block until signal received
		<-sigint
		cancelFn()
	}()

	s.Start(ctx)
}

func (f flags) validate() error {
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
