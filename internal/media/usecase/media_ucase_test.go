package usecase_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/media"
	"source.toby3d.me/toby3d/pub/internal/media/usecase"
)

func TestUpload(t *testing.T) {
	t.Parallel()

	f := domain.TestFile(t)
	repo := media.NewSpyMediaRepository(media.NewStubMediaRepository(f, nil))

	out, err := usecase.NewMediaUseCase(repo).Upload(context.Background(), *f)
	if err != nil {
		t.Fatal(err)
	}

	if out.Path == "" {
		t.Error("expect non-empty location path, got nothing")
	}

	if expect := 1; repo.Creates != expect {
		t.Errorf("expect %d Create calls, got %d", expect, repo.Creates)
	}
}

func TestDownload(t *testing.T) {
	t.Parallel()

	f := domain.TestFile(t)
	repo := media.NewStubMediaRepository(f, nil)

	out, err := usecase.NewMediaUseCase(repo).
		Download(context.Background(), f.Path)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(f, out); diff != "" {
		t.Errorf("%#+v", diff)
	}
}
