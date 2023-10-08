package usecase

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/entry"
)

type entryUseCase struct {
	entries entry.Repository
}

// Create implements entry.UseCase.
func (ucase *entryUseCase) Create(ctx context.Context, e domain.Entry) (*domain.Entry, error) {
	result, err := ucase.entries.Create(ctx, e.URL.RequestURI(), e)
	if err != nil {
		return nil, fmt.Errorf("cannot source entry: %w", err)
	}

	return result, nil
}

// Delete implements entry.UseCase.
func (ucase *entryUseCase) Delete(ctx context.Context, u *url.URL) (bool, error) {
	if _, err := ucase.entries.Update(ctx, u.RequestURI(), func(_ context.Context, e *domain.Entry) (
		*domain.Entry, error,
	) {
		now := time.Now().UTC()
		e.DeletedAt = now
		e.UpdatedAt = now

		return nil, nil
	}); err != nil {
		return false, fmt.Errorf("cannot undelete entry: %w", err)
	}

	return true, nil
}

// Source implements entry.UseCase.
func (ucase *entryUseCase) Source(ctx context.Context, u *url.URL) (*domain.Entry, error) {
	result, err := ucase.entries.Get(ctx, u.RequestURI())
	if err != nil {
		return nil, fmt.Errorf("cannot source entry: %w", err)
	}

	return result, nil
}

// Undelete implements entry.UseCase.
func (ucase *entryUseCase) Undelete(ctx context.Context, u *url.URL) (*domain.Entry, error) {
	result, err := ucase.entries.Update(ctx, u.RequestURI(), func(_ context.Context, e *domain.Entry) (
		*domain.Entry, error,
	) {
		e.DeletedAt = time.Time{}
		e.UpdatedAt = time.Now().UTC()

		return nil, nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot undelete entry: %w", err)
	}

	return result, nil
}

// Update implements entry.UseCase.
func (ucase *entryUseCase) Update(ctx context.Context, u *url.URL, e domain.Entry) (*domain.Entry, error) {
	result, err := ucase.entries.Update(ctx, u.RequestURI(), func(_ context.Context, e *domain.Entry) (
		*domain.Entry, error,
	) {
		e.DeletedAt = time.Time{}
		e.UpdatedAt = time.Now().UTC()

		return nil, nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot undelete entry: %w", err)
	}

	return result, nil
}

func NewEntryUseCase(entries entry.Repository) entry.UseCase {
	return &entryUseCase{
		entries: entries,
	}
}
