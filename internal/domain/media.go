package domain

import (
	_ "embed"
	"testing"
)

type Media struct {
	Name        string
	ContentType string
	Content     []byte
}

//go:embed testdata/sunset.jpg
var testMediaContent []byte

func TestMedia(tb testing.TB) *Media {
	tb.Helper()

	return &Media{
		Name:        "sunset.jpg",
		ContentType: "image/jpeg",
		Content:     testMediaContent,
	}
}
