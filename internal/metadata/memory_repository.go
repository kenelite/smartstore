package metadata

import (
	"context"
	"sync"
	"time"
)

// InMemoryRepository is useful for local dev / fallback when DB is not configured.
type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*ObjectRecord
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*ObjectRecord),
	}
}

func makeKey(env, region, bucket, key string) string {
	return env + "|" + region + "|" + bucket + "|" + key
}

func (r *InMemoryRepository) GetObject(_ context.Context, env, region, bucket, key string) (*ObjectRecord, error) {
	k := makeKey(env, region, bucket, key)
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.data[k]
	if !ok || rec.Status == "DELETED" {
		return nil, ErrNotFound
	}
	return rec, nil
}

func (r *InMemoryRepository) PutObject(_ context.Context, rec *ObjectRecord) error {
	if rec == nil {
		return ErrNotFound
	}
	now := time.Now()
	rec.UpdatedAt = now
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = now
	}
	if rec.Version == 0 {
		rec.Version = 1
	}
	if rec.Status == "" {
		rec.Status = "ACTIVE"
	}
	k := makeKey(rec.Env, rec.LogicalRegion, rec.Bucket, rec.ObjectKey)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[k] = rec
	return nil
}

func (r *InMemoryRepository) MarkDeleted(_ context.Context, env, region, bucket, key string) error {
	k := makeKey(env, region, bucket, key)
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.data[k]
	if !ok {
		return ErrNotFound
	}
	rec.Status = "DELETED"
	rec.UpdatedAt = time.Now()
	return nil
}
