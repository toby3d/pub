package memory

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/entry"
)

type memoryEntryRepository struct {
	mutex   *sync.RWMutex
	entries map[string]domain.Entry
}

func NewMemoryEntryRepository() entry.Repository {
	return &memoryEntryRepository{
		mutex:   new(sync.RWMutex),
		entries: make(map[string]domain.Entry),
	}
}

func (repo *memoryEntryRepository) Create(ctx context.Context, p string, e domain.Entry) error {
	p = path.Clean(strings.ToLower(p))

	_, err := repo.Get(ctx, p)
	if err != nil && !errors.Is(err, entry.ErrNotExist) {
		return fmt.Errorf("cannot save a new entry: %w", err)
	}
	if err == nil {
		return entry.ErrExist
	}

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.entries[p] = e

	return nil
}

func (repo *memoryEntryRepository) Get(ctx context.Context, p string) (*domain.Entry, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	if out, ok := repo.entries[path.Clean(strings.ToLower(p))]; ok {
		return &out, nil
	}

	return nil, entry.ErrNotExist
}

func (repo *memoryEntryRepository) Fetch(ctx context.Context, p string) ([]domain.Entry, int, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	out := make([]domain.Entry, 0)

	for entryPath, e := range repo.entries {
		if matched, err := path.Match(p+"*", entryPath); err != nil || !matched {
			continue
		}

		out = append(out, e)
	}

	return out, len(out), nil
}

func (repo *memoryEntryRepository) Update(ctx context.Context, p string, update entry.UpdateFunc) (*domain.Entry, error) {
	p = path.Clean(strings.ToLower(p))

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if in, ok := repo.entries[p]; ok {
		out, err := update(ctx, &in)
		if err != nil {
			return nil, fmt.Errorf("cannot update entry: %w", err)
		}

		repo.entries[p] = *out

		return out, err
	}

	return nil, fmt.Errorf("cannot update entry: %w", entry.ErrNotExist)
}

func (repo *memoryEntryRepository) Delete(ctx context.Context, p string) (bool, error) {
	if _, err := repo.Get(ctx, p); err != nil {
		if errors.Is(err, entry.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("cannot find entry to delete: %w", err)
	}

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	delete(repo.entries, path.Clean(strings.ToLower(p)))

	return true, nil
}
