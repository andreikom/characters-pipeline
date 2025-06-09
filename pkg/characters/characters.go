package characters

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"charactersync/model"
	"charactersync/providers"
	"charactersync/store"
)

type Service struct {
	fetchers []providers.Fetcher
	store    store.Storer
}

func NewCharactersService(fetchers []providers.Fetcher, store store.Storer) Service {
	return Service{
		fetchers: fetchers,
		store:    store,
	}
}

func (s Service) FetchAll(ctx context.Context, opts *providers.FetchOptions) ([]*model.Character, error) {
	var all []*model.Character

	for _, f := range s.fetchers {
		chars, err := f.Fetch(ctx, opts)
		if err != nil {
			fmt.Printf("warning: provider %T failed: %v\n", f, err)
			continue
		}
		all = append(all, chars...)
	}

	sort.Slice(all, func(i, j int) bool {
		return strings.ToLower(all[i].Name) < strings.ToLower(all[j].Name)
	})

	if err := s.store.Persist(all); err != nil {
		return nil, fmt.Errorf("failed to persist characters: %w", err)
	}

	return all, nil
}
