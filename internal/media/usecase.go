package media

import (
	"context"
	"net/url"

	"source.toby3d.me/toby3d/pub/internal/domain"
)

type (
	UseCase interface {
		// Upload uploads media file into micropub store which can be
		// download later.
		Upload(ctx context.Context, file domain.File) (*url.URL, error)

		// Download downloads early uploaded media stored in path.
		Download(ctx context.Context, path string) (*domain.File, error)
	}

	dummyUseCase struct{}

	stubUseCase struct {
		u    *url.URL
		file *domain.File
		err  error
	}
)

// NewDummyUseCase creates a dummy use case what does nothing.
func NewDummyUseCase() UseCase {
	return &dummyUseCase{}
}

func (dummyUseCase) Upload(_ context.Context, _ domain.File) (*url.URL, error)  { return nil, nil }
func (dummyUseCase) Download(_ context.Context, _ string) (*domain.File, error) { return nil, nil }

// NewDummyUseCase creates a stub use case what always returns provided input.
func NewStubUseCase(err error, file *domain.File, u *url.URL) UseCase {
	return &stubUseCase{
		u:    u,
		file: file,
		err:  err,
	}
}

func (ucase stubUseCase) Upload(_ context.Context, _ domain.File) (*url.URL, error) {
	return ucase.u, ucase.err
}

func (ucase stubUseCase) Download(_ context.Context, _ string) (*domain.File, error) {
	return ucase.file, ucase.err
}
