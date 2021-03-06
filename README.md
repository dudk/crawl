## Recursive crawler

Crawler prints crawled URLs to stdout.

#### Input

* URL

#### Design

The crawler is built on top of the following astractions:

```go
type (
	// Scheduler manages fetching execution.
	Scheduler interface {
		Start(context.Context, Fetcher) chan struct{}
	}

	// Fetcher is responsible for fetching content of the web page and
	// calling BodyReader and .
	Fetcher interface {
		Fetch(context.Context, Visitor, BodyReader, string) []string
	}

	// BodyReader reads the HTTP response body. It must consume all data
	// available in io.Reader.
	BodyReader interface {
		ReadBody(context.Context, io.Reader) error
	}

	// BodyReader reads the HTTP response body. It must consume all data
	// available in io.Reader. BodyReader must respect context and return
	// if it's cancelled.
	BodyReader interface {
		ReadBody(context.Context, *url.URL, io.Reader) error
	}

	// Visitor records if parsed URL was visited. It's safe for concurrent
	// use.
	Visitor interface {
		Visit(string) bool
	}
)
```

### Instructions

Run the following commands:

```shell
go build
./crawl -url https://godoc.org/
```
