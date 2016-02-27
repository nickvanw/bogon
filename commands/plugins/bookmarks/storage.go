package bookmarks

import "sync"

// Storage is an interface to store and retrieve bookmarks in k/b form
type Storage interface {
	New(string, string) error
	Lookup(string) (string, bool)
	Remove(string) error
	Dump() (map[string]string, error)
}

// MemoryStorage maps bookmark k/v pairs in memory
// It is thread-safe
type MemoryStorage struct {
	sync.RWMutex
	d map[string]string
}

// NewMemStorage returns a ready to use memory bookmark storage
func NewMemStorage() *MemoryStorage {
	return &MemoryStorage{d: map[string]string{}}
}

// New creates a new bookmark at key with the contents of value
func (m *MemoryStorage) New(key, value string) error {
	m.Lock()
	m.d[key] = value
	m.Unlock()
	return nil
}

// Lookup attempts to find the bookmark at the key, returning true if it is there
func (m *MemoryStorage) Lookup(key string) (string, bool) {
	m.RLock()
	data, ok := m.d[key]
	m.RUnlock()
	return data, ok
}

// Remove deletes a key if it's in the datastore
func (m *MemoryStorage) Remove(key string) error {
	m.Lock()
	delete(m.d, key)
	m.Unlock()
	return nil
}

func (m *MemoryStorage) Dump() (map[string]string, error) {
	m.RLock()
	defer m.RUnlock()
	return m.d, nil
}
