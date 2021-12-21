package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"source.toby3d.me/website/micropub/internal/media"
)

type mediaUseCase struct {
	repo media.Repository
}

const DefaultNameLength = 64

func NewMediaUseCase(repo media.Repository) media.UseCase {
	return &mediaUseCase{
		repo: repo,
	}
}

func (useCase *mediaUseCase) Upload(ctx context.Context, name string, src []byte) (string, error) {
	newName := make([]byte, DefaultNameLength)
	if _, err := rand.Read(newName); err != nil {
		return "", fmt.Errorf("cannot generate random string: %w", err)
	}

	fileName := base64.RawURLEncoding.EncodeToString(newName) + filepath.Ext(name)
	if err := useCase.repo.Create(ctx, fileName, src); err != nil {
		return "", fmt.Errorf("cannot create media: %w", err)
	}

	return fileName, nil
}

func (useCase *mediaUseCase) Download(ctx context.Context, name string) ([]byte, error) {
	result, err := useCase.repo.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("cannot find media: %w", err)
	}

	return result, nil
}
