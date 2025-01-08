package vault

import (
	"errors"
	"fmt"
)

type Entity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Vault struct {
	storage Storable
}

var ErrCommit = errors.New("failed to commit to the storage")

func New(config Config) *Vault {
	storage := newStorage(config)

	err := storage.Init()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize the storage: %v", err))
	}

	return &Vault{storage: storage}
}

func (v *Vault) All() Records {
	return v.storage.All()
}

func (v *Vault) Get(key string) string {
	return v.storage.Get(key)
}

func (v *Vault) Upsert(key, value string) error {
	if err := v.storage.Upsert(key, value); err != nil {
		return fmt.Errorf("failed to upsert: %w", err)
	}

	return nil
}

func (v *Vault) Destroy(key string) error {
	if err := v.storage.Destroy(key); err != nil {
		return fmt.Errorf("failed to destroy: %w", err)
	}

	return nil
}
