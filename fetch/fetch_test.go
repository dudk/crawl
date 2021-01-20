package fetch_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
			f.Parser = fetch.ParseBodyFunc(func(_ *url.URL, r io.Reader) ([]string, error) {
				return nil, fmt.Errorf("parser error")
			})
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
		Parser: fetch.ParseBodyFunc(func(_ *url.URL, r io.Reader) ([]string, error) {
			s := bufio.NewScanner(r)
			var result []string
			for s.Scan() {
				result = append(result, s.Text())
			}
			return result, nil
		}),
	}

	t.Run("fetcher ok", testOk(f))
	t.Run("context done", testContextDone(f))
	t.Run("parser error", testParserError(f))

}

type bodyReader bytes.Buffer

func (br *bodyReader) ReadBody(ctx context.Context, URL *url.URL, r io.Reader) error {
	_, err := io.Copy((*bytes.Buffer)(br), r)
	return err
}

func (br *bodyReader) body() string {
	return string((*bytes.Buffer)(br).Bytes())
}

func TestBodyReader(t *testing.T) {
	testBodyOk := func(f fetch.Fetcher) func(*testing.T) {
		return func(t *testing.T) {
			var br bodyReader
			f.BodyReader = &br
			bodyURLs := []string{"http://b", "http://c"}
			response := fmt.Sprintf("%s\n%s\n", bodyURLs[0], bodyURLs[1])
			ts := httptest.NewServer(
				http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					res.Write([]byte(response))
				}),
			)
			defer ts.Close()
			urls, err := f.Fetch(context.Background(), ts.URL)
			assertEqual(t, "fetch error", err, nil)
			assertEqual(t, "urls", urls, bodyURLs)
			assertEqual(t, "body", br.body(), response)
		}
	}
	testBodyParserErr := func(f fetch.Fetcher) func(*testing.T) {
		return func(t *testing.T) {
			f.Parser = fetch.ParseBodyFunc(func(_ *url.URL, r io.Reader) ([]string, error) {
				return nil, fmt.Errorf("parser error")
			})
			var br bodyReader
			f.BodyReader = &br
			bodyURLs := []string{"http://b", "http://c"}
			response := fmt.Sprintf("%s\n%s\n", bodyURLs[0], bodyURLs[1])
			ts := httptest.NewServer(
				http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					res.Write([]byte(response))
				}),
			)
			defer ts.Close()
			urls, err := f.Fetch(context.Background(), ts.URL)
			assertEqual(t, "fetch error", strings.Contains(err.Error(), "parser error"), true)
			assertEqual(t, "urls", urls, []string(nil))
			assertEqual(t, "body", br.body(), "")
		}
	}
	testBodyReaderErr := func(f fetch.Fetcher) func(*testing.T) {
		return func(t *testing.T) {
			f.BodyReader = fetch.BodyReaderFunc(func(context.Context, *url.URL, io.Reader) error {
				return fmt.Errorf("reader error")
			})
			bodyURLs := []string{"http://b", "http://c"}
			response := fmt.Sprintf("%s\n%s\n", bodyURLs[0], bodyURLs[1])
			ts := httptest.NewServer(
				http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					res.Write([]byte(response))
				}),
			)
			defer ts.Close()
			urls, err := f.Fetch(context.Background(), ts.URL)
			assertEqual(t, "fetch error", strings.Contains(err.Error(), "reader error"), true)
			assertEqual(t, "urls", urls, []string(nil))
		}
	}
	testParserAndBodyReaderErr := func(f fetch.Fetcher) func(*testing.T) {
		return func(t *testing.T) {
			f.Parser = fetch.ParseBodyFunc(func(_ *url.URL, r io.Reader) ([]string, error) {
				return nil, fmt.Errorf("parser error")
			})
			f.BodyReader = fetch.BodyReaderFunc(func(context.Context, *url.URL, io.Reader) error {
				return fmt.Errorf("reader error")
			})
			bodyURLs := []string{"http://b", "http://c"}
			response := fmt.Sprintf("%s\n%s\n", bodyURLs[0], bodyURLs[1])
			ts := httptest.NewServer(
				http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					res.Write([]byte(response))
				}),
			)
			defer ts.Close()
			urls, err := f.Fetch(context.Background(), ts.URL)
			assertEqual(t, "fetch error", strings.Contains(err.Error(), "reader error"), true)
			assertEqual(t, "fetch error", strings.Contains(err.Error(), "parser error"), true)
			assertEqual(t, "urls", urls, []string(nil))
		}
	}

	f := fetch.Fetcher{
		Client: http.DefaultClient,
		Parser: fetch.ParseBodyFunc(func(_ *url.URL, r io.Reader) ([]string, error) {
			s := bufio.NewScanner(r)
			var result []string
			for s.Scan() {
				result = append(result, s.Text())
			}
			return result, nil
		}),
	}

	t.Run("body reader ok", testBodyOk(f))
	t.Run("body parser error", testBodyParserErr(f))
	t.Run("body reader error", testBodyReaderErr(f))
	t.Run("body reader error", testParserAndBodyReaderErr(f))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}
