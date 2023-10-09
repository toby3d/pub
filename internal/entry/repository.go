package entry

import (
	"context"
	"errors"

	"source.toby3d.me/toby3d/pub/internal/domain"
)

type (
	UpdateFunc func(ctx context.Context, input *domain.Entry) (*domain.Entry, error)

	Repository interface {
		Create(ctx context.Context, path string, e domain.Entry) error
		Get(ctx context.Context, path string) (*domain.Entry, error)
		Fetch(ctx context.Context, path string) ([]domain.Entry, int, error)
		Update(ctx context.Context, path string, update UpdateFunc) (*domain.Entry, error)
		Delete(ctx context.Context, path string) (bool, error)
	}
)

var (
	ErrExist    error = errors.New("this entry already exist")
	ErrNotExist error = errors.New("this entry is not exist")
)
