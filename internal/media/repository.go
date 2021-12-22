package media

import (
	"context"
	"errors"

	"source.toby3d.me/website/micropub/internal/domain"
)

type Repository interface {
	// Create save media contents by provided name.
	Create(ctx context.Context, name string, media *domain.Media) error

	// Delete remove early saved media file contents by name.
	Delete(ctx context.Context, name string) error

	// Get returns media file contents by name.
	Get(ctx context.Context, name string) (*domain.Media, error)
}

var ErrNotExist = errors.New("media not exist")
