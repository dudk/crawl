package fetch_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/dudk/crawl/fetch"
)

func TestFetcher(t *testing.T) {
	testOk := func(f fetch.Fetcher) func(*testing.T) {
		return func(t *testing.T) {
			bodyURLs := []string{"http://b", "http://c"}
			ts := httptest.NewServer(
				http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					fmt.Fprintf(res, "%s\n%s\n", bodyURLs[0], bodyURLs[1])
				}),
			)
			defer ts.Close()
			urls, err := f.Fetch(context.Background(), ts.URL)
			assertEqual(t, "urls", urls, bodyURLs)
			assertEqual(t, "fetch error", err, nil)
		}
	}
	testContextDone := func(f fetch.Fetcher) func(*testing.T) {
		return func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			cancelFn()
			urls, err := f.Fetch(ctx, "http://localhost")
			assertEqual(t, "urls", len(urls), 0)
			assertEqual(t, "fetch error", strings.Contains(err.Error(), "context canceled"), true)
		}
	}
	testParserError := func(f fetch.Fetcher) func(*testing.T) {
		return func(t *testing.T) {
			f.Parser = func(r io.Reader) ([]string, error) {
				return nil, fmt.Errorf("parser error")
			}
			bodyURLs := []string{"http://b", "http://c"}
			ts := httptest.NewServer(
				http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					fmt.Fprintf(res, "%s\n%s\n", bodyURLs[0], bodyURLs[1])
				}),
			)
			defer ts.Close()
			urls, err := f.Fetch(context.Background(), ts.URL)
			assertEqual(t, "urls", len(urls), 0)
			assertEqual(t, "fetch error", strings.Contains(err.Error(), "parser error"), true)
		}
	}
	f := fetch.Fetcher{
		Client: http.DefaultClient,
		Parser: func(r io.Reader) ([]string, error) {
			s := bufio.NewScanner(r)
			var result []string
			for s.Scan() {
				result = append(result, s.Text())
			}
			return result, nil
		},
	}

	t.Run("fetcher ok", testOk(f))
	t.Run("context done", testContextDone(f))
	t.Run("parser error", testParserError(f))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}
