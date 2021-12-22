package media

import (
	"context"

	"source.toby3d.me/website/micropub/internal/domain"
)

type UseCase interface {
	// Upload save uploaded media file in temporary storage with random
	// generated name.
	Upload(ctx context.Context, media *domain.Media) (string, error)

	// Download returns early uploaded media file by random generated name.
	Download(ctx context.Context, name string) (*domain.Media, error)
}
