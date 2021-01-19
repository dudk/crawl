package parsing_test

import (
	"bytes"
	"net/url"
	"reflect"
	"testing"

	"github.com/dudk/crawl/caching"
	"github.com/dudk/crawl/parsing"
)

func TestParser(t *testing.T) {
	p := parsing.HTML{
		BaseURL: "http://base.org",
		Visitor: caching.NewInMemoryCache(),
	}
	expected := []string{
		"http://base.org/a",
		"http://base.org/b",
		"http://base.org/c",
		"http://base.org/c/d",
		"http://base.org/c/f",
	}

	buf := bytes.NewBuffer([]byte(`<html>
	<head></head>
	<body>
		<a href="http://base.org/a"></a>
		<a href="http://base.org/b"></a>
		<a href="http://base.org/c"></a>
		<a href="http://base.org/c/d"></a>
		<a href="http://base.net/e"></a>
		<a href="f"></a>
	</body>
	</html>`))
	u, _ := url.Parse("http://base.org/c/")
	links, _ := p.ParseBody(u, buf)
	assertEqual(t, "urls", links, expected)
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}
