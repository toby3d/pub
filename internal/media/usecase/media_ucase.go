package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"

	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/media"
)

type mediaUseCase struct {
	media media.Repository
}

const charset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func NewMediaUseCase(media media.Repository) media.UseCase {
	return &mediaUseCase{
		media: media,
	}
}

func (ucase *mediaUseCase) Upload(ctx context.Context, file domain.File) (*url.URL, error) {
	randName := make([]byte, 64)

	for i := range randName {
		randName[i] = charset[rand.Intn(len(charset))]
	}

	newName := string(randName) + "." + file.Ext()

	if err := ucase.media.Create(ctx, newName, file); err != nil {
		return nil, fmt.Errorf("cannot upload nedia: %w", err)
	}

	return &url.URL{
		Path: newName,
	}, nil
}

func (ucase *mediaUseCase) Download(ctx context.Context, path string) (*domain.File, error) {
	out, err := ucase.media.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("cannot find media file: %w", err)
	}

	return out, nil
}
