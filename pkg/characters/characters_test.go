package characters_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"charactersync/model"
	"charactersync/pkg/characters"
	"charactersync/providers"
)

type mockFetcher struct {
	characters []*model.Character
	err        error
}

func (d *mockFetcher) Fetch(ctx context.Context, opts *providers.FetchOptions) ([]*model.Character, error) {
	return d.characters, d.err
}

type mockStore struct {
	saved []*model.Character
	err   error
}

func (ds *mockStore) Persist(chars []*model.Character) error {
	if ds.err != nil {
		return ds.err
	}
	ds.saved = chars
	return nil
}

func TestService_FetchAll_Success(t *testing.T) {
	mockChars1 := []*model.Character{
		{Name: "Bulbasaur", Origin: "Pokemon", Species: "Grass", AdditionalAttribute: "64"},
		{Name: "Luke Skywalker", Origin: "Star Wars", Species: "Human", AdditionalAttribute: "19BBY"},
	}
	mockChars2 := []*model.Character{
		{Name: "Bulbasaur", Origin: "Pokemon", Species: "Grass", AdditionalAttribute: "64"},
		{Name: "Luke Skywalker", Origin: "Star Wars", Species: "Human", AdditionalAttribute: "19BBY"},
	}

	service := characters.NewCharactersService(
		[]providers.Fetcher{
			&mockFetcher{characters: mockChars1},
			&mockFetcher{characters: mockChars2},
		},
		&mockStore{},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := service.FetchAll(ctx, &providers.FetchOptions{LimitPerProvider: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != (len(mockChars1) + len(mockChars2)) {
		t.Fatalf("expected %d characters, got %d", len(mockChars1)+len(mockChars2), len(result))
	}
}

func TestService_FetchAll_FetcherError(t *testing.T) {
	service := characters.NewCharactersService(
		[]providers.Fetcher{
			&mockFetcher{err: errors.New("fetcher failed")},
		},
		&mockStore{},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := service.FetchAll(ctx, &providers.FetchOptions{})
	if err != nil {
		t.Fatalf("should not return error even when fetcher fails, got: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result on fetch failure, got: %d", len(result))
	}
}

func TestService_FetchAll_PersistError(t *testing.T) {
	mockChars := []*model.Character{
		{Name: "Morty", Origin: "Rick and Morty", Species: "Human", AdditionalAttribute: "Alive"},
	}
	service := characters.NewCharactersService(
		[]providers.Fetcher{
			&mockFetcher{characters: mockChars},
		},
		&mockStore{err: errors.New("db failed")},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := service.FetchAll(ctx, &providers.FetchOptions{})
	if err == nil {
		t.Fatal("expected error from Persist, got nil")
	}
}

type limitTrackingFetcher struct {
	expectedLimit int
	called        bool
	t             *testing.T
}

func (f *limitTrackingFetcher) Fetch(ctx context.Context, opts *providers.FetchOptions) ([]*model.Character, error) {
	f.called = true
	if opts == nil {
		f.t.Fatalf("expected FetchOptions, got nil")
	}
	if opts.LimitPerProvider != f.expectedLimit {
		f.t.Errorf("expected limit %d, got %d", f.expectedLimit, opts.LimitPerProvider)
	}

	return []*model.Character{
		{Name: "Test", Origin: "TestOrigin", Species: "TestSpecies", AdditionalAttribute: "123"},
	}, nil
}

func TestService_FetchAll_RespectsLimit(t *testing.T) {
	limit := 5
	limitFetcher := &limitTrackingFetcher{
		expectedLimit: limit,
		t:             t,
	}

	service := characters.NewCharactersService(
		[]providers.Fetcher{limitFetcher},
		&mockStore{},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := service.FetchAll(ctx, &providers.FetchOptions{LimitPerProvider: limit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !limitFetcher.called {
		t.Error("expected fetcher to be called")
	}
}

func TestFetchAll_SortedByName(t *testing.T) {
	mockCharacters := []*model.Character{
		{Name: "Luke Skywalker", Origin: "Star Wars", Species: "Human", AdditionalAttribute: "19BBY"},
		{Name: "Leia Organa", Origin: "Star Wars", Species: "Human", AdditionalAttribute: "19BBY"},
		{Name: "Han Solo", Origin: "Star Wars", Species: "Human", AdditionalAttribute: "29BBY"},
	}

	service := characters.NewCharactersService(
		[]providers.Fetcher{
			&mockFetcher{characters: mockCharacters},
		},
		&mockStore{},
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	opts := &providers.FetchOptions{LimitPerProvider: 10}

	result, err := service.FetchAll(ctx, opts)
	if err != nil {
		t.Fatalf("FetchAll failed: %v", err)
	}

	for i := 1; i < len(result); i++ {
		if result[i-1].Name > result[i].Name {
			t.Errorf("Characters are not sorted by name: %v > %v", result[i-1].Name, result[i].Name)
		}
	}
}
