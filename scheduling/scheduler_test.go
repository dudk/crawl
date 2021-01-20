package scheduling_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/dudk/crawl/caching"
	"github.com/dudk/crawl/fetch"
	"github.com/dudk/crawl/parsing"
	"github.com/dudk/crawl/scheduling"
)

func TestScheduler(t *testing.T) {
	baseURL := "https://godoc.org/"
	v := caching.NewInMemoryCache()
	s := scheduling.UnboundScheduler{
		BaseURL: baseURL,
		ErrorHandler: scheduling.ErrorHandlerFunc(func(err error) {
			fmt.Fprintf(os.Stdout, "fetched error: %s\n", err.Error())
		}),
		Fetcher: fetch.Fetcher{
			Client: http.DefaultClient,
			Parser: parsing.HTML{
				BaseURL: baseURL,
				Visitor: v,
			},
			BodyReader: fetch.BodyReaderFunc(func(ctx context.Context, URL *url.URL, r io.Reader) error {
				fmt.Fprintf(os.Stdout, "fetched URL: %s\n", URL.String())
				if _, err := io.Copy(ioutil.Discard, r); err != nil {
					return fmt.Errorf("error discarding body: %w", err)
				}
				return nil
			}),
		},
	}
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*1)
	defer cancelFn()
	s.Start(ctx)
}
