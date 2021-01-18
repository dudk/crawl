package caching_test

import (
	"reflect"
	"testing"

	"github.com/dudk/crawl/caching"
)

func TestInMemoryCache(t *testing.T) {
	c := caching.NewInMemoryCache()
	s := "value"
	ok := c.Visit(s)
	assertEqual(t, "first visit", ok, false)

	ok = c.Visit(s)
	assertEqual(t, "second visit", ok, true)
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}
