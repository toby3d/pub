package media

import (
	"context"
	"errors"

	"source.toby3d.me/toby3d/pub/internal/domain"
)

type (
	UpdateFunc func(src *domain.File) (*domain.File, error)

	Repository interface {
		// Create save provided file into the store as a new media.
		// Returns error if media already exists.
		Create(ctx context.Context, path string, file domain.File) error

		// Get returns a early stored media as a file. Returns error if
		// file is not exist.
		Get(ctx context.Context, path string) (*domain.File, error)

		// Update replaces already exists media file or creates a new
		// one if it is not. Returns error overwise.
		Update(ctx context.Context, path string, update UpdateFunc) error

		// Delete removes media file from the store. Returns error if
		// existed file cannot be or already deleted.
		Delete(ctx context.Context, path string) error
	}

	dummyRepository struct{}

	stubRepository struct {
		output *domain.File
		err    error
	}

	spyRepository struct {
		subRepository Repository
		Calls         int
		Creates       int
		Updates       int
		Gets          int
		Deletes       int
	}

	// NOTE(toby3d): fakeRepository is already provided by memory sub-package.
	// NOTE(toby3d): mockRepository is complicated. Mocking too much is bad.
)

var (
	ErrExist    error = errors.New("this file already exist")
	ErrNotExist error = errors.New("this file is not exist")
)

// NewDummyMediaRepository creates an empty repository to satisfy contracts.
// It is used in tests where repository working is not important.
func NewDummyMediaRepository() Repository {
	return &dummyRepository{}
}

func (dummyRepository) Create(_ context.Context, _ string, _ domain.File) error { return nil }
func (dummyRepository) Get(_ context.Context, _ string) (*domain.File, error)   { return nil, nil }
func (dummyRepository) Update(_ context.Context, _ string, _ UpdateFunc) error  { return nil }
func (dummyRepository) Delete(_ context.Context, _ string) error                { return nil }

// NewStubMediaRepository creates a repository that always returns input as a
// output. It is used in tests where some dependency on the repository is
// required.
func NewStubMediaRepository(output *domain.File, err error) Repository {
	return &stubRepository{
		output: output,
		err:    err,
	}
}

func (repo *stubRepository) Create(_ context.Context, _ string, _ domain.File) error {
	return repo.err
}

func (repo *stubRepository) Get(_ context.Context, _ string) (*domain.File, error) {
	return repo.output, repo.err
}

func (repo *stubRepository) Update(_ context.Context, _ string, _ UpdateFunc) error {
	return repo.err
}

func (repo *stubRepository) Delete(_ context.Context, _ string) error {
	return repo.err
}

// NewSpyMediaRepository creates a spy repository which count outside calls,
// based on provided subRepo. If subRepo is nil, then DummyRepository will be
// used.
func NewSpyMediaRepository(subRepo Repository) *spyRepository {
	if subRepo == nil {
		subRepo = NewDummyMediaRepository()
	}

	return &spyRepository{
		subRepository: subRepo,
		Creates:       0,
		Updates:       0,
		Gets:          0,
		Deletes:       0,
	}
}

func (repo *spyRepository) Create(_ context.Context, _ string, _ domain.File) error {
	repo.Creates++

	return repo.subRepository.Create(context.TODO(), "", domain.File{})
}

func (repo *spyRepository) Get(_ context.Context, _ string) (*domain.File, error) {
	repo.Gets++

	return repo.subRepository.Get(context.TODO(), "")
}

func (repo *spyRepository) Update(_ context.Context, _ string, _ UpdateFunc) error {
	repo.Updates++

	return repo.subRepository.Update(context.TODO(), "", nil)
}

func (repo *spyRepository) Delete(_ context.Context, _ string) error {
	repo.Deletes++

	return repo.subRepository.Delete(context.TODO(), "")
}
