package entry

import (
	"context"
	"errors"

	"source.toby3d.me/toby3d/pub/internal/domain"
)

type (
	UpdateFunc func(ctx context.Context, input *domain.Entry) (*domain.Entry, error)

	Repository interface {
		Create(ctx context.Context, path string, e domain.Entry) error
		Get(ctx context.Context, path string) (*domain.Entry, error)
		Fetch(ctx context.Context, path string) ([]domain.Entry, int, error)
		Update(ctx context.Context, path string, update UpdateFunc) (*domain.Entry, error)
		Delete(ctx context.Context, path string) (bool, error)
	}

	dummyRepository struct{}

	stubRepository struct {
		outputs []domain.Entry
		output  *domain.Entry
		err     error
		ok      bool
	}

	spyRepository struct {
		subRepository Repository
		Calls         int
		Creates       int
		Deletes       int
		Fetches       int
		Gets          int
		Updates       int
	}

	// NOTE(toby3d): fakeRepository is already provided by memory sub-package.
	// NOTE(toby3d): mockRepository is complicated. Mocking too much is bad.
)

var (
	ErrExist    error = errors.New("this entry already exist")
	ErrNotExist error = errors.New("this entry is not exist")
)

// NewDummyMediaRepository creates an empty repository to satisfy contracts.
// It is used in tests where repository working is not important.
func NewDummyEntryRepository() Repository {
	return &dummyRepository{}
}

func (dummyRepository) Create(_ context.Context, _ string, _ domain.Entry) error { return nil }
func (dummyRepository) Delete(_ context.Context, _ string) (bool, error)         { return false, nil }
func (dummyRepository) Get(_ context.Context, _ string) (*domain.Entry, error)   { return nil, nil }

func (dummyRepository) Fetch(_ context.Context, _ string) ([]domain.Entry, int, error) {
	return make([]domain.Entry, 0), 0, nil
}

func (dummyRepository) Update(_ context.Context, _ string, _ UpdateFunc) (*domain.Entry, error) {
	return nil, nil
}

// NewStubEntryRepository creates a repository that always returns input as a
// output. It is used in tests where some dependency on the repository is
// required.
func NewStubEntryRepository(outputs []domain.Entry, output *domain.Entry, err error, ok bool) Repository {
	return &stubRepository{
		outputs: outputs,
		output:  output,
		err:     err,
		ok:      ok,
	}
}

func (repo *stubRepository) Create(_ context.Context, _ string, _ domain.Entry) error {
	return repo.err
}

func (repo *stubRepository) Delete(ctx context.Context, path string) (bool, error) {
	return repo.ok, repo.err
}

func (repo *stubRepository) Fetch(ctx context.Context, path string) ([]domain.Entry, int, error) {
	return repo.outputs, len(repo.outputs), repo.err
}

func (repo *stubRepository) Get(ctx context.Context, path string) (*domain.Entry, error) {
	return repo.output, repo.err
}

func (repo *stubRepository) Update(ctx context.Context, path string, update UpdateFunc) (*domain.Entry, error) {
	return repo.output, repo.err
}

// NewSpyEntryRepository creates a spy repository which count outside calls,
// based on provided subRepo. If subRepo is nil, then DummyRepository will be
// used.
func NewSpyEntryRepository(subRepo Repository) *spyRepository {
	if subRepo == nil {
		subRepo = NewDummyEntryRepository()
	}

	return &spyRepository{
		subRepository: subRepo,
		Creates:       0,
		Updates:       0,
		Gets:          0,
		Fetches:       0,
		Deletes:       0,
	}
}

func (repo *spyRepository) Create(ctx context.Context, path string, e domain.Entry) error {
	repo.Creates++

	return repo.subRepository.Create(ctx, path, e)
}

func (repo *spyRepository) Delete(ctx context.Context, path string) (bool, error) {
	repo.Deletes++

	return repo.subRepository.Delete(ctx, path)
}

func (repo *spyRepository) Fetch(ctx context.Context, path string) ([]domain.Entry, int, error) {
	repo.Fetches++

	return repo.subRepository.Fetch(ctx, path)
}

func (repo *spyRepository) Get(ctx context.Context, path string) (*domain.Entry, error) {
	repo.Gets++

	return repo.subRepository.Get(ctx, path)
}

func (repo *spyRepository) Update(ctx context.Context, path string, update UpdateFunc) (*domain.Entry, error) {
	repo.Updates++

	return repo.subRepository.Update(ctx, path, update)
}
