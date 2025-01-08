package vault

import "path/filepath"

type Records map[string]string

type Storable interface {
	Init() error
	All() Records
	Exists(key string) bool
	Get(key string) string
	Upsert(key, value string) error
	Destroy(key string) error
}

const FILE_STORAGE = "file"
const MEMORY_STORAGE = "memory"
const FILE_STORAGE_NAME = "storage.enc"

func newStorage(config Config) Storable {
	switch config.Storage {
	case FILE_STORAGE:
		filename := filepath.Join(projectRoot(), FILE_STORAGE_NAME)
		return &FileStorage{filename: filename, secret: config.Secret}
	default:
		return &MemoryStorage{}
	}
}
