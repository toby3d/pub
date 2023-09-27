package domain

import (
	"net/url"
	"path"
	"testing"
	"time"
)

// Entry represent a single microformats2 entry.
type Entry struct {
	UpdatedAt    time.Time
	PublishedAt  time.Time
	DeletedAt    time.Time
	URL          *url.URL
	Params       map[string]any
	Title        string
	Description  string
	Photo        []*url.URL
	Syndications []*url.URL
	Content      []byte
	Tags         []string
}

// TestEntry returns a random valid Entry for tests.
func TestEntry(tb testing.TB) *Entry {
	tb.Helper()

	return &Entry{
		URL:   &url.URL{Path: path.Join("/", "samples", "lipsum")},
		Title: "Lorem ipsum",
		Description: "Ut enim ad minim veniam, quis nostrud exercitation " +
			"ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		Content: []byte("Duis aute irure dolor in reprehenderit in " +
			"voluptate velit esse cillum dolore eu fugiat nulla " +
			"pariatur. Excepteur sint occaecat cupidatat non proident," +
			" sut in culpa qui officia deserunt mollit anim id est " +
			"laborum."),
		Tags: []string{"lorem", "ipsum", "dor"},
	}
}
