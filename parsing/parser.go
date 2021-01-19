package parsing

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Visitor records if parsed URL was visited. It's safe for concurrent
// use.
type Visitor interface {
	Visit(string) bool
}

// HTML is dead-stupid html parser. It uses standard library url package to
// normalize URLs. Only URLs having base URL as prefix will be returned by the parser.
type HTML struct {
	BaseURL string
	Visitor
}

// ParseBody parses HTML for URLS. Only <a> tags are examined and exact URL
// matching is made.
func (p HTML) ParseBody(url *url.URL, r io.Reader) ([]string, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("HTML parsing error: %w", err)
	}

	var links []string
	return parseLinks(links, p.Visitor, n, url, p.BaseURL), nil
}

// recursive depth-first html parsing.
func parseLinks(links []string, v Visitor, n *html.Node, url *url.URL, baseURL string) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			if URL, ok := normalize(url, a.Val); ok && strings.HasPrefix(URL, baseURL) && !v.Visit(URL) {
				links = append(links, URL)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = parseLinks(links, v, c, url, baseURL)
	}
	return links
}

// validate and normalize raw URL. returns false if URL cannot be parsed.
func normalize(baseURL *url.URL, URL string) (string, bool) {
	u, err := url.Parse(URL)
	if err != nil {
		return URL, false
	}
	u.Fragment = ""
	if u.IsAbs() {
		return u.String(), true
	}

	u = baseURL.ResolveReference(u)
	return u.String(), true
}
