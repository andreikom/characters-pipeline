package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"charactersync/model"
)

const swapiBaseURL = "https://swapi.py4e.com/api"

type SWAPIFetcher struct {
	client  *http.Client
	baseURL string
}

func NewSWAPIFetcher() *SWAPIFetcher {
	return &SWAPIFetcher{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: swapiBaseURL,
	}
}

type swapiPeopleResponse struct {
	Next    string        `json:"next"`
	Results []swapiPerson `json:"results"`
}

type swapiPerson struct {
	Name      string   `json:"name"`
	BirthYear string   `json:"birth_year"`
	Species   []string `json:"species"`
}

type swapiSpecies struct {
	Name string `json:"name"`
}

func (sf *SWAPIFetcher) Fetch(ctx context.Context, opts *FetchOptions) ([]*model.Character, error) {
	var starwarsCharacters []*model.Character
	url := fmt.Sprintf("%s/people", sf.baseURL)
	fetched := 0

	for url != "" && (opts.LimitPerProvider == 0 || fetched < opts.LimitPerProvider) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := sf.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var page swapiPeopleResponse
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			return nil, err
		}

		for _, person := range page.Results {
			if opts.LimitPerProvider > 0 && fetched >= opts.LimitPerProvider {
				break
			}

			speciesName := "Unknown"
			if len(person.Species) > 0 {
				speciesReq, err := http.NewRequestWithContext(ctx, http.MethodGet, person.Species[0], nil)
				if err == nil {
					speciesResp, err := sf.client.Do(speciesReq)
					if err == nil {
						defer speciesResp.Body.Close()
						var species swapiSpecies
						if err := json.NewDecoder(speciesResp.Body).Decode(&species); err == nil {
							speciesName = species.Name
						}
					}
				}
			}

			starwarsCharacters = append(starwarsCharacters, &model.Character{
				Name:                person.Name,
				Origin:              "Star Wars",
				Species:             speciesName,
				AdditionalAttribute: person.BirthYear,
			})
			fetched++
		}

		url = page.Next
	}

	return starwarsCharacters, nil
}
