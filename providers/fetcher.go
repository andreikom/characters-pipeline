package providers

import (
	"context"

	"charactersync/model"
)

type Fetcher interface {
	Fetch(ctx context.Context, opts *FetchOptions) ([]*model.Character, error)
}

type FetchOptions struct {
	LimitPerProvider int
}
