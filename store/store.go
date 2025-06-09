package store

import (
	"charactersync/model"
)

type Storer interface {
	Persist(characters []*model.Character) error
}
