package entry

import (
	"context"
	"net/url"

	"source.toby3d.me/toby3d/pub/internal/domain"
)

type (
	UseCase interface {
		// Create creates a new entry. Returns map or rel links, like Permalink
		// or created post, shortcode and syndication.
		Create(ctx context.Context, e domain.Entry) (map[string]*url.URL, error)

		// Update updates exist entry properties on provided u.
		//
		// TODO(toby3d): return Location header if entry updates their URL.
		Update(ctx context.Context, u *url.URL, e domain.Entry) (*domain.Entry, error)

		// Delete destroy entry on provided URL.
		Delete(ctx context.Context, u *url.URL) (bool, error)

		// Undelete recover deleted entry on provided URL.
		Undelete(ctx context.Context, u *url.URL) (*domain.Entry, error)

		// Source returns properties of entry on provided URL.
		Source(ctx context.Context, u *url.URL) (*domain.Entry, error)
	}

	stubUseCase struct{}
)

func NewStubUseCase() *stubUseCase {
	return &stubUseCase{}
}

func (ucase *stubUseCase) Create(ctx context.Context, e domain.Entry) (map[string]*url.URL, error) {
	return nil, nil
}

func (ucase *stubUseCase) Update(ctx context.Context, u *url.URL, e domain.Entry) (*domain.Entry, error) {
	return nil, nil
}

func (ucase *stubUseCase) Delete(ctx context.Context, u *url.URL) (bool, error) { return false, nil }

func (ucase *stubUseCase) Undelete(ctx context.Context, u *url.URL) (*domain.Entry, error) {
	return nil, nil
}

func (ucase *stubUseCase) Source(ctx context.Context, u *url.URL) (*domain.Entry, error) {
	return nil, nil
}
