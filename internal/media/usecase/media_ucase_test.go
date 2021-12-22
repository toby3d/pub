package usecase_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"source.toby3d.me/website/micropub/internal/domain"
	repository "source.toby3d.me/website/micropub/internal/media/repository/memory"
	"source.toby3d.me/website/micropub/internal/media/usecase"
)

func TestUpload(t *testing.T) {
	t.Parallel()

	media := domain.TestMedia(t)

	result, err := usecase.NewMediaUseCase(repository.NewMemoryMediaRepository(new(sync.Map))).
		Upload(context.Background(), media)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestDownload(t *testing.T) {
	t.Parallel()

	media := domain.TestMedia(t)

	fileName := make([]byte, usecase.DefaultNameLength)
	_, err := rand.Read(fileName)
	require.NoError(t, err)

	newName := base64.RawURLEncoding.EncodeToString(fileName) + filepath.Ext(media.Name)

	repo := repository.NewMemoryMediaRepository(new(sync.Map))
	require.NoError(t, repo.Create(context.Background(), newName, media))

	result, err := usecase.NewMediaUseCase(repo).Download(context.Background(), newName)
	assert.NoError(t, err)
	assert.Equal(t, result, media)
}
