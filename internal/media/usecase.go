package media

import (
	"context"
)

type UseCase interface {
	// Upload save uploaded media file in temporary storage with random
	// generated name.
	Upload(ctx context.Context, name string, contents []byte) (string, error)

	// Download returns early uploaded media file by random generated name.
	Download(ctx context.Context, name string) ([]byte, error)
}
