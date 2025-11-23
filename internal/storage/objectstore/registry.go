package objectstore

import "sync"

// ProviderRegistry maps provider name to ObjectStorage implementation.
type ProviderRegistry struct {
	mu       sync.RWMutex
	backends map[string]ObjectStorage
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		backends: make(map[string]ObjectStorage),
	}
}

func (r *ProviderRegistry) Register(name string, backend ObjectStorage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends[name] = backend
}

func (r *ProviderRegistry) Get(name string) (ObjectStorage, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.backends[name]
	return b, ok
}
