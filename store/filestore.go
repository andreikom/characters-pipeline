package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"charactersync/model"
)

type FileStore struct {
	filepathBase string
}

func NewFileStore(filepathBase string) *FileStore {
	return &FileStore{
		filepathBase: filepathBase,
	}
}

func (fs *FileStore) Persist(characters []*model.Character) error {
	timestamp := time.Now().Format("20250102_120102")
	filename := fmt.Sprintf("%s%s_%s.json", fs.filepathBase, "characters", timestamp)
	fullpath := filepath.Clean(filename)

	file, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fullpath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(characters); err != nil {
		return fmt.Errorf("failed to encode characters: %w", err)
	}

	return nil
}
