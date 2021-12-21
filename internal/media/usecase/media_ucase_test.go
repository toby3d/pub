package usecase_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repository "source.toby3d.me/website/micropub/internal/media/repository/memory"
	"source.toby3d.me/website/micropub/internal/media/usecase"
)

func TestUpload(t *testing.T) {
	t.Parallel()

	_, contents := testFile(t)

	result, err := usecase.NewMediaUseCase(repository.NewMemoryMediaRepository(new(sync.Map))).
		Upload(context.Background(), "sunset.jpg", contents)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestDownload(t *testing.T) {
	t.Parallel()

	fileName, contents := testFile(t)

	repo := repository.NewMemoryMediaRepository(new(sync.Map))
	require.NoError(t, repo.Create(context.Background(), fileName, contents))

	result, err := usecase.NewMediaUseCase(repo).Download(context.Background(), fileName)
	assert.NoError(t, err)
	assert.Equal(t, result, contents)
}

func testFile(tb testing.TB) (string, []byte) {
	tb.Helper()

	fileName := make([]byte, usecase.DefaultNameLength)
	_, err := rand.Read(fileName)
	require.NoError(tb, err)

	contents := make([]byte, 128)
	_, err = rand.Read(contents)
	require.NoError(tb, err)

	return base64.RawURLEncoding.EncodeToString(fileName) + ".jpg", contents
}
