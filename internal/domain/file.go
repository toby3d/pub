package domain

import (
	"embed"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"testing"
)

// File represent a single media file, like photo.
type File struct {
	Path    string // content/example/photo.jpg
	Content []byte
}

//go:embed testdata/sunset.jpg
var testdata embed.FS

// TestFile returns a valid File for tests.
func TestFile(tb testing.TB) *File {
	tb.Helper()

	f, err := testdata.Open(filepath.Join("testdata", "sunset.jpg"))
	if err != nil {
		tb.Fatalf("cannot open testing file: %s", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		tb.Fatalf("cannot fetch testing file info: %s", err)
	}

	body, err := io.ReadAll(f)
	if err != nil {
		tb.Fatalf("cannot read testing file body: %s", err)
	}

	return &File{
		Path:    info.Name(),
		Content: body,
	}
}

// LogicalName returns full file name without directory path.
func (f File) LogicalName() string {
	return filepath.Base(f.Path)
}

// BaseFileName returns file name without extention and directory path.
func (f File) BaseFileName() string {
	base := filepath.Base(f.Path)

	return strings.TrimSuffix(base, filepath.Ext(base))
}

// Ext returns file extention.
func (f File) Ext() string {
	return filepath.Ext(f.Path)
}

// Dir returns file directory.
func (f File) Dir() string {
	return filepath.Dir(f.Path)
}

// MediaType returns media type based on file extention.
func (f File) MediaType() string {
	return mime.TypeByExtension(f.Ext())
}
