## Recursive crawler

#### Input

* URL
* Path

#### Design

The crawler is built on top of the following astractions:

```go
type (
	// Scheduler manages fetching execution.
	Scheduler interface {
		Start(context.Context, Fetcher) chan struct{}
	}

	// Fetcher is responsible for fetching content of the web page and
	// calling BodyReader.
	Fetcher interface {
		Fetch(context.Context, Visitor, BodyReader, string) []string
	}

	// BodyReader reads the HTTP response body. It must consume all data
	// available in io.Reader.
	BodyReader interface {
		ReadBody(context.Context, io.Reader) error
	}

	// Visitor records if parsed URL was visited. It's safe for concurrent
	// use.
	Visitor interface {
		Visit(string) bool
	}
)
```