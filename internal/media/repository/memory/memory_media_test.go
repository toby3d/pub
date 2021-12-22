package memory_test

import (
	"context"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"source.toby3d.me/website/micropub/internal/domain"
	repository "source.toby3d.me/website/micropub/internal/media/repository/memory"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	media := domain.TestMedia(t)

	store := new(sync.Map)
	require.NoError(t, repository.NewMemoryMediaRepository(store).
		Create(context.Background(), "sample.ext", media))

	result, ok := store.Load(path.Join(repository.DefaultPathPrefix, "sample.ext"))
	assert.True(t, ok)
	assert.Equal(t, media, result)
}

func TestGet(t *testing.T) {
	t.Parallel()

	media := domain.TestMedia(t)

	store := new(sync.Map)
	store.Store(path.Join(repository.DefaultPathPrefix, "sample.ext"), media)

	result, err := repository.NewMemoryMediaRepository(store).Get(context.Background(), "sample.ext")
	assert.NoError(t, err)
	assert.Equal(t, media, result)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	media := domain.TestMedia(t)

	store := new(sync.Map)
	store.Store(path.Join(repository.DefaultPathPrefix, "sample.ext"), media)

	require.NoError(t, repository.NewMemoryMediaRepository(store).Delete(context.Background(), "sample.ext"))

	result, ok := store.Load(path.Join(repository.DefaultPathPrefix, "sample.ext"))
	assert.False(t, ok)
	assert.Nil(t, result)
}
