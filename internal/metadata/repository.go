package metadata

import (
	"context"
	"errors"
	"time"
)

type StoreBackend string

const (
	StoreRedisOnly   StoreBackend = "REDIS_ONLY"
	StoreObjectOnly  StoreBackend = "OBJECT_ONLY"
	StoreRedisObject StoreBackend = "REDIS_OBJECT"
)

type ObjectRecord struct {
	Env           string
	LogicalRegion string
	Bucket        string
	ObjectKey     string

	SizeBytes    int64
	ContentType  string
	StorageClass string
	StoreBackend StoreBackend

	ProviderType   string
	ProviderRegion string
	ProviderBucket string
	PhysicalKey    string

	ETag    string
	Version int64
	Status  string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	GetObject(ctx context.Context, env, region, bucket, key string) (*ObjectRecord, error)
	PutObject(ctx context.Context, rec *ObjectRecord) error
	MarkDeleted(ctx context.Context, env, region, bucket, key string) error
}

var ErrNotFound = errors.New("object not found")
