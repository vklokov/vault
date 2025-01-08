package vault

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func projectRoot() string {
	ex, err := os.Executable()
	if err != nil {
		panic("unable to access root dir")
	}

	return filepath.Dir(ex)
}

var mu sync.Mutex

type FileStorage struct {
	secret      string
	filename    string
	initialized bool
	records     Records
}

func (s *FileStorage) Init() error {
	s.records = make(Records)

	f, err := os.OpenFile(s.filename, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if len(body) != 0 {
		data, err := decryptSrting(string(body), s.secret)
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(data), &s.records)
		if err != nil {
			return err
		}
	}

	s.initialized = true

	return nil
}

func (s *FileStorage) Exists(key string) bool {
	_, ok := s.records[key]
	return ok
}

func (s *FileStorage) All() Records {
	s.checkStorageReady()
	return s.records
}

func (s *FileStorage) Get(key string) string {
	s.checkStorageReady()

	return s.records[key]
}

func (s *FileStorage) Upsert(key, value string) error {
	s.checkStorageReady()
	s.records[key] = value

	if err := s.rewrite(); err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) Destroy(key string) error {
	s.checkStorageReady()
	delete(s.records, key)

	if err := s.rewrite(); err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) checkStorageReady() {
	if !s.initialized {
		panic("storage must be initialized")
	}
}

func (s *FileStorage) rewrite() error {
	data, err := json.Marshal(s.records)
	if err != nil {
		return err
	}

	body, err := encryptString(string(data), s.secret)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(s.filename, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	f.Write([]byte(body))

	return nil
}
