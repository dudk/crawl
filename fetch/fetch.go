package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	// BodyParser parses the HTTP response body and returns a list of URLs
	// found in the body.
	BodyParser interface {
		ParseBody(URL *url.URL, r io.Reader) ([]string, error)
	}

	// ParseBodyFunc implements BodyParser.
	ParseBodyFunc func(URL *url.URL, r io.Reader) ([]string, error)

	// BodyReader reads the HTTP response body. It must consume all data
	// available in io.Reader. BodyReader must respect context and return
	// if it's cancelled.
	BodyReader interface {
		ReadBody(context.Context, *url.URL, io.Reader) error
	}

	// BodyReaderFunc implements BodyReader.
	BodyReaderFunc func(context.Context, *url.URL, io.Reader) error
)

// Fetcher is responsible for recursive fetch of pages.
type Fetcher struct {
	Client     *http.Client
	Parser     BodyParser
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

	if f.BodyReader == nil {
		return f.parseOnly(res)
	}
	return f.parseAndRead(ctx, res)
}

func (f Fetcher) parseOnly(res *http.Response) ([]string, error) {
	links, err := f.Parser.ParseBody(res.Request.URL, res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing links: %w", err)
	}
	return links, nil
}

func (f Fetcher) parseAndRead(ctx context.Context, res *http.Response) ([]string, error) {
	pr, pw := io.Pipe()
	// tee reader writes everything from HTTP body to pipe
	tr := io.TeeReader(res.Body, pw)

	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		// close pipe once we are done
		defer pw.Close()

		// read HTTP body and write to pipe
		err := f.BodyReader.ReadBody(ctx, res.Request.URL, tr)
		if err != nil {
			errChan <- err
		}
	}()

	// read body from pipe
	var err error
	links, err := f.Parser.ParseBody(res.Request.URL, pr)
	if err != nil {
		// if pipe reader has failed, we need to close writer
		pw.Close()
		err = fmt.Errorf("parsing error: %w", err)
	}

	readErr := <-errChan
	if readErr != nil && readErr != io.ErrClosedPipe {
		if err != nil {
			return nil, fmt.Errorf("%s and body read error: %w", err.Error(), readErr)
		}
		return nil, fmt.Errorf("body read error: %w", readErr)
	}

	return links, err
}

// ReadBody implements BodyReader.
func (fn BodyReaderFunc) ReadBody(ctx context.Context, URL *url.URL, r io.Reader) error {
	return fn(ctx, URL, r)
}

// ParseBody implements BodyParser.
func (fn ParseBodyFunc) ParseBody(u *url.URL, r io.Reader) ([]string, error) {
	return fn(u, r)
}
