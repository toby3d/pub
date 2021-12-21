package media

import (
	"context"
	"errors"
)

type Repository interface {
	// Create save media contents by provided name.
	Create(ctx context.Context, name string, contents []byte) error

	// Delete remove early saved media file contents by name.
	Delete(ctx context.Context, name string) error

	// Get returns media file contents by name.
	Get(ctx context.Context, name string) ([]byte, error)
}

var ErrNotFound = errors.New("media not found")
