package memory

import (
	"context"
	"path/filepath"
	"sync"

	"source.toby3d.me/website/micropub/internal/media"
)

type memoryMediaRepository struct {
	store *sync.Map
}

const DefaultPathPrefix string = "media"

func NewMemoryMediaRepository(store *sync.Map) media.Repository {
	return &memoryMediaRepository{
		store: store,
	}
}

func (repo *memoryMediaRepository) Create(ctx context.Context, name string, contents []byte) error {
	repo.store.Store(filepath.Join(DefaultPathPrefix, name), contents)

	return nil
}

func (repo *memoryMediaRepository) Get(ctx context.Context, name string) ([]byte, error) {
	src, ok := repo.store.Load(filepath.Join(DefaultPathPrefix, name))
	if !ok {
		return nil, media.ErrNotFound
	}

	contents, ok := src.([]byte)
	if !ok {
		return nil, media.ErrNotFound
	}

	return contents, nil
}

func (repo *memoryMediaRepository) Delete(ctx context.Context, name string) error {
	repo.store.Delete(filepath.Join(DefaultPathPrefix, name))

	return nil
}
