package vault

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
const FILE_STORAGE_NAME = "storage.enc"
const MEMORY_STORAGE = "memory"

func newStorage(config Config) Storable {
	switch config.Storage {
	case FILE_STORAGE:
		return &FileStorage{filename: FILE_STORAGE_NAME, secret: config.Secret}
	default:
		return &MemoryStorage{}
	}
}
