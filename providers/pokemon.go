package providers

import (
	"context"
	"fmt"
	"strconv"

	pokeapi "github.com/JoshGuarino/PokeGo/pkg"

	"charactersync/model"
)

// PokemonFetcher - note: "https://pokeapi.co/api/v2/" baseURL is embedded in the PokeGo library
type PokemonFetcher struct {
	client pokeapi.PokeGo
}

func NewPokemonFetcher() *PokemonFetcher {
	return &PokemonFetcher{
		client: pokeapi.NewClient(),
	}
}

func (pf *PokemonFetcher) Fetch(ctx context.Context, opts *FetchOptions) ([]*model.Character, error) {
	var pokemonCharacters []*model.Character
	limit := opts.LimitPerProvider
	if limit <= 0 {
		limit = 100 // default page size
	}

	offset := 0
	for {
		page, err := pf.client.Pokemon.GetPokemonList(limit, offset)
		if err != nil {
			return nil, fmt.Errorf("PokeGo list error: %w", err)
		}

		for _, entry := range page.Results {
			p, err := pf.client.Pokemon.GetPokemon(entry.Name)
			if err != nil {
				// skip individual failures
				continue
			}

			pokemonCharacters = append(pokemonCharacters, &model.Character{
				Name:                p.Name,
				Origin:              "Pokemon",
				Species:             p.Species.Name,
				AdditionalAttribute: strconv.Itoa(p.BaseExperience),
			})

			// stop if we've hit the per provider limit
			if len(pokemonCharacters) >= opts.LimitPerProvider && opts.LimitPerProvider > 0 {
				return pokemonCharacters[:opts.LimitPerProvider], nil
			}
		}

		// if fewer than we requested, we're at the end
		if len(page.Results) < limit {
			break
		}
		offset += limit
	}

	return pokemonCharacters, nil
}
