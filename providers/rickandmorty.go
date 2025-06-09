package providers

import (
	"context"
	"fmt"

	"charactersync/model"
	rickandmorty "github.com/pitakill/rickandmortyapigowrapper"
)

// RickAndMortyFetcher - note: "https://rickandmortyapi.com/api/" baseURL is embedded in the rickandmortyapigowrapper library
type RickAndMortyFetcher struct{}

func NewRickAndMortyFetcher() *RickAndMortyFetcher {
	return &RickAndMortyFetcher{}
}

func (rf *RickAndMortyFetcher) Fetch(ctx context.Context, opts *FetchOptions) ([]*model.Character, error) {
	var rockAndMortyCharacters []*model.Character
	page := 1
	fetched := 0
	limit := opts.LimitPerProvider

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		options := map[string]interface{}{
			"endpoint": "character",
			"page":     page,
		}

		data, err := rickandmorty.GetCharacters(options)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		for _, c := range data.Results {
			if limit > 0 && fetched >= limit {
				return rockAndMortyCharacters, nil
			}

			char := &model.Character{
				Name:                c.Name,
				Origin:              "Rick and Morty",
				Species:             c.Species,
				AdditionalAttribute: c.Status,
			}

			rockAndMortyCharacters = append(rockAndMortyCharacters, char)
			fetched++
		}

		if data.Info.Next == "" || (limit > 0 && fetched >= limit) {
			break
		}

		page++
	}

	return rockAndMortyCharacters, nil
}
