package memory_test

import (
	"context"
	"crypto/rand"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repository "source.toby3d.me/website/micropub/internal/media/repository/memory"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	fileName, contents := testFile(t)
	store := new(sync.Map)

	require.NoError(t, repository.NewMemoryMediaRepository(store).
		Create(context.Background(), fileName, contents))

	result, ok := store.Load(path.Join(repository.DefaultPathPrefix, fileName))
	assert.True(t, ok)
	assert.Equal(t, result, contents)
}

func TestGet(t *testing.T) {
	t.Parallel()

	fileName, contents := testFile(t)
	store := new(sync.Map)

	store.Store(path.Join(repository.DefaultPathPrefix, fileName), contents)

	result, err := repository.NewMemoryMediaRepository(store).Get(context.Background(), fileName)
	assert.NoError(t, err)
	assert.Equal(t, result, contents)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	fileName, contents := testFile(t)
	store := new(sync.Map)

	store.Store(path.Join(repository.DefaultPathPrefix, fileName), contents)

	require.NoError(t, repository.NewMemoryMediaRepository(store).Delete(context.Background(), fileName))

	result, ok := store.Load(path.Join(repository.DefaultPathPrefix, fileName))
	assert.False(t, ok)
	assert.Nil(t, result)
}

func testFile(tb testing.TB) (string, []byte) {
	tb.Helper()

	contents := make([]byte, 128)
	_, err := rand.Read(contents)
	require.NoError(tb, err)

	return "sunset.jpg", contents
}
