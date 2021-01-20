package scheduling

import (
	"context"
	"sync"
)

type (
	// UnboundScheduler schedules recursively URL fetching.
	UnboundScheduler struct {
		Fetcher      Fetcher
		ErrorHandler ErrorHandler
		BaseURL      string
		wg           sync.WaitGroup
	}

	// Fetcher is responsible for fetching content of the web page and
	// calling BodyReader.
	Fetcher interface {
		Fetch(ctx context.Context, url string) ([]string, error)
	}

	// ErrorHandler handles fetch error.
	ErrorHandler interface {
		Handle(err error)
	}

	// ErrorHandlerFunc implements ErrorHandler.
	ErrorHandlerFunc func(err error)
)

// Handle the occured error.
func (fn ErrorHandlerFunc) Handle(err error) {
	fn(err)
}

// Start scheduler execution. This call blocks goroutine until all URLs are
// fetched or context is cancelled.
func (s *UnboundScheduler) Start(ctx context.Context) {
	s.wg.Add(1)
	go s.fetch(ctx, s.BaseURL)
	s.wg.Wait()
}

func (s *UnboundScheduler) fetch(ctx context.Context, url string) {
	urls, err := s.Fetcher.Fetch(ctx, url)
	if err != nil && s.ErrorHandler != nil {
		s.ErrorHandler.Handle(err)
	}
	for _, u := range urls {
		s.wg.Add(1)
		go s.fetch(ctx, u)
	}
	s.wg.Done()
}
