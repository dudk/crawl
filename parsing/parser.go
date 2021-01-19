package parsing

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// HTML is dead-stupid html parser. It uses standard library url package to
// normalize URLs. Only URLs having base URL as prefix will be returned by the parser.
type HTML struct {
	BaseURL string
}

// ParseBody parses HTML for URLS. Only <a> tags are examined and exact URL
// matching is made.
func (p HTML) ParseBody(r io.Reader) ([]string, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("HTML parsing error: %w", err)
	}

	var links []string
	return parseLinks(links, n, p.BaseURL), nil
}

// recursive depth-first html parsing.
func parseLinks(links []string, n *html.Node, baseURL string) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			if URL, ok := normalize(a.Val); ok && strings.HasPrefix(URL, baseURL) {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = parseLinks(links, c, baseURL)
	}
	return links
}

// validate and normalize raw URL. returns false if URL cannot be parsed.
func normalize(URL string) (string, bool) {
	u, err := url.Parse(URL)
	if err != nil {
		return URL, false
	}
	return u.String(), true
}
