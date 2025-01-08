package vault

type MemoryStorage struct {
	records Records
}

func (m *MemoryStorage) Init() error {
	return nil
}

func (m *MemoryStorage) All() Records {
	return m.records
}

func (m *MemoryStorage) Exists(key string) bool {
	_, ok := m.records[key]
	return ok
}

func (m *MemoryStorage) Get(key string) string {
	return m.records[key]
}

func (m *MemoryStorage) Upsert(key, value string) error {
	m.records[key] = value
	return nil
}

func (m *MemoryStorage) Destroy(key string) error {
	delete(m.records, key)
	return nil
}
