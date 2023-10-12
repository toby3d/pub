package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"source.toby3d.me/toby3d/pub/internal/domain"
	"source.toby3d.me/toby3d/pub/internal/entry"
	"source.toby3d.me/toby3d/pub/internal/entry/usecase"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	e := domain.TestEntry(t)
	repo := entry.NewSpyEntryRepository(entry.NewStubEntryRepository(nil, e, nil, false))

	if _, err := usecase.NewEntryUseCase(repo).Create(context.Background(), *e); err != nil {
		t.Fatal(err)
	}

	if repo.Creates == 0 {
		t.Error("expect creation call")
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	e := domain.TestEntry(t)

	for name, tc := range map[string]struct {
		options entry.UpdateOptions
		expect  func() *domain.Entry
	}{
		"add": {
			options: entry.UpdateOptions{
				Add: &domain.Entry{Tags: []string{"indieweb", "testing"}},
			},
			expect: func() *domain.Entry {
				updated := *e
				updated.Tags = append(updated.Tags, "indieweb", "testing")

				return &updated
			},
		},
		// TODO(toby3d): "update": {},
		// TODO(toby3d): "delete": {},
	} {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			expect := tc.expect()
			repo := entry.NewSpyEntryRepository(entry.NewStubEntryRepository(nil, expect, nil, false))
			ucase := usecase.NewEntryUseCase(repo)

			out, err := ucase.Update(context.Background(), e.URL, tc.options)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(out, expect, cmp.AllowUnexported(e.RSVP)); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	e := domain.TestEntry(t)
	deleted := *e
	deleted.DeletedAt = time.Now().UTC()

	repo := entry.NewSpyEntryRepository(entry.NewStubEntryRepository(nil, &deleted, nil, false))

	ok, err := usecase.NewEntryUseCase(repo).Delete(context.Background(), e.URL)
	if err != nil {
		t.Fatal(err)
	}

	if !ok || repo.Updates == 0 || repo.Deletes != 0 {
		t.Errorf("expect update call without deleting")
	}
}

func TestUndelete(t *testing.T) {
	t.Parallel()

	e := domain.TestEntry(t)
	undeleted := *e
	undeleted.DeletedAt = time.Now().UTC().AddDate(0, 0, -7)

	repo := entry.NewSpyEntryRepository(entry.NewStubEntryRepository(nil, e, nil, false))

	if _, err := usecase.NewEntryUseCase(repo).Undelete(context.Background(), e.URL); err != nil {
		t.Fatal(err)
	}

	if repo.Updates == 0 {
		t.Error("expect update call")
	}
}

func TestSource(t *testing.T) {
	t.Parallel()

	e := domain.TestEntry(t)
	repo := entry.NewSpyEntryRepository(entry.NewStubEntryRepository(nil, e, nil, false))

	if _, err := usecase.NewEntryUseCase(repo).Source(context.Background(), e.URL); err != nil {
		t.Fatal(err)
	}

	if repo.Gets == 0 {
		t.Error("expect getting call")
	}
}
