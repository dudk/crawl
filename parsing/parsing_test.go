package parsing_test

import (
	"bytes"
	"testing"

	"github.com/dudk/crawl/parsing"
)

func TestParser(t *testing.T) {
	p := parsing.HTML{
		BaseURL: "http://base.org",
	}
	results := []string{
		"http://base.org/a",
		"http://base.org/b",
		"http://base.org/c",
		"http://base.org/c/d",
	}

	buf := bytes.NewBuffer([]byte(`<html>
	<head></head>
	<body>
		<a href="http://base.org/a"></a>
		<a href="http://base.org/b"></a>
		<a href="http://base.org/c"></a>
		<a href="http://base.org/c/d"></a>
		<a href="http://base.net/e"></a>
		<a href="https://base.org/f"></a>
	</body>
	</html>`))

	links, _ := p.ParseBody(buf)
	if len(links) != len(results) {
		t.Fatalf("expected: %v got: %v", results, links)
	}

}
