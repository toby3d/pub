package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"source.toby3d.me/website/micropub/internal/domain"
	"source.toby3d.me/website/micropub/internal/media"
)

type mediaUseCase struct {
	repo media.Repository
}

const DefaultNameLength = 32

func NewMediaUseCase(repo media.Repository) media.UseCase {
	return &mediaUseCase{
		repo: repo,
	}
}

func (useCase *mediaUseCase) Upload(ctx context.Context, media *domain.Media) (string, error) {
	newName := make([]byte, DefaultNameLength)
	if _, err := rand.Read(newName); err != nil {
		return "", fmt.Errorf("cannot generate random string: %w", err)
	}

	fileName := base64.RawURLEncoding.EncodeToString(newName) + filepath.Ext(media.Name)
	if err := useCase.repo.Create(ctx, fileName, media); err != nil {
		return "", fmt.Errorf("cannot create media: %w", err)
	}

	return fileName, nil
}

func (useCase *mediaUseCase) Download(ctx context.Context, name string) (*domain.Media, error) {
	result, err := useCase.repo.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("cannot find media: %w", err)
	}

	return result, nil
}
