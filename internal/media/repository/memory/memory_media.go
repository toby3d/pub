package memory

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/media"
)

type memoryMediaRepository struct {
	mutex *sync.RWMutex
	media map[string]domain.File
}

func NewMemoryMediaRepository() media.Repository {
	return &memoryMediaRepository{
		mutex: new(sync.RWMutex),
		media: make(map[string]domain.File),
	}
}

func (repo *memoryMediaRepository) Create(ctx context.Context, p string, f domain.File) error {
	p = path.Clean(strings.ToLower(p))

	_, err := repo.Get(ctx, p)
	if err != nil && !errors.Is(err, media.ErrNotExist) {
		return fmt.Errorf("cannot save a new media: %w", err)
	}
	if err == nil {
		return media.ErrExist
	}

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.media[p] = f

	return nil
}

func (repo *memoryMediaRepository) Get(ctx context.Context, p string) (*domain.File, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	if out, ok := repo.media[path.Clean(strings.ToLower(p))]; ok {
		return &out, nil
	}

	return nil, media.ErrNotExist
}

func (repo *memoryMediaRepository) Update(ctx context.Context, p string, update media.UpdateFunc) error {
	p = path.Clean(strings.ToLower(p))

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if in, ok := repo.media[p]; ok {
		out, err := update(&in)
		if err != nil {
			return fmt.Errorf("cannot update media: %w", err)
		}

		repo.media[p] = *out

		return nil
	}

	return media.ErrNotExist
}

func (repo *memoryMediaRepository) Delete(ctx context.Context, p string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	delete(repo.media, path.Clean(strings.ToLower(p)))

	return nil
}
