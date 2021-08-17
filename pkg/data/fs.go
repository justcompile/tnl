package data

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

type FSStore struct {
	config map[string]string
}

func (f *FSStore) Get(_ context.Context, key string) (string, error) {
	if val, exists := f.config[key]; exists {
		return val, nil
	}
	return "", errors.New("not found")
}

func (f *FSStore) Save(context.Context, string, string) error {
	return nil
}

func NewFSStore(filePath string) (*FSStore, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	store := &FSStore{}
	if err := json.NewDecoder(f).Decode(&store.config); err != nil {
		return nil, err
	}

	return store, nil
}
