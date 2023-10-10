package domain

import (
	"net/url"
	"path"
	"testing"
	"time"
)

// Entry represent a single microformats2 entry.
type Entry struct {
	CreatedAt   time.Time
	DeletedAt   time.Time
	Title       string    // p-name
	Description string    // p-summary
	Content     []byte    // e-content
	PublishedAt time.Time // dt-published
	UpdatedAt   time.Time // dt-updated
	// TODO(toby3d): Author string // p-author
	Tags []string // p-category
	URL  *url.URL // u-url
	ID   string   // u-uid
	// TODO(toby3d): Location string // p-location
	Syndications []*url.URL // u-syndication
	// TODO(toby3d): Reply *url.URL // u-in-reply-to
	RSVP RSVP // p-rsvp
	// TODO(toby3d): Like *url.URL // u-like-of
	// TODO(toby3d): Repost *url.URL // u-repost-of

	// Draft Properties
	// TODO(toby3d): Comments []string // p-comment
	Photo []*url.URL // u-photo
	Video []*url.URL // u-video

	// Proposed Additions
	Audio []*url.URL // u-audio
	// TODO(toby3d): Like []*url.URL // u-like
	// TODO(toby3d): Repost []*url.URL // u-repost
	// TODO(toby3d): BookmarkOf []*url.URL // u-bookmark-of
	// TODO(toby3d): Featured []*url.URL // u-featured
	Latitude  float32 // p-latitude
	Longitude float32 // p-longitude
	Altitude  float32 // p-altitude
	// TODO(toby3d): Duration int // p-duration
	// TODO(toby3d): Size int // p-size
	// TODO(toby3d): ListenOf *url.URL // u-listen-of
	// TODO(toby3d): WatchOf *url.URL // u-watch-of
	// TODO(toby3d): ReadOf *url.URL // u-read-of
	// TODO(toby3d): TranslationOf *url.URL // u-translation-of
	// TODO(toby3d): Checkin *url.URL // u-checking
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
