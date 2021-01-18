package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type (
	ParseFunc func(r io.Reader) ([]string, error)

	// BodyReader reads the HTTP response body. It must consume all data
	// available in io.Reader.
	BodyReader interface {
		ReadBody(context.Context, io.Reader) error
	}
)

// Fetcher is responsible for recursive fetch of pages.
type Fetcher struct {
	BaseURL    string
	Client     *http.Client
	Parser     ParseFunc
	BodyReader BodyReader
}

// Fetch takes in context and returns a list of parsed urls.
func (f Fetcher) Fetch(ctx context.Context, url string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := f.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer res.Body.Close()

	links, err := f.Parser(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing links: %w", err)
	}
	return links, nil
}
