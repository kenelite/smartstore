package smart

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/kenelite/smartstore/internal/cache"
	"github.com/kenelite/smartstore/internal/metadata"
	"github.com/kenelite/smartstore/internal/storage/objectstore"
)

type Service struct {
	cache     *cache.RedisCache
	metaRepo  metadata.Repository
	router    objectstore.ObjectRoute
	providers *objectstore.ProviderRegistry

	smallFileThreshold int64         // bytes, e.g. 1MB
	cacheTTL           time.Duration // TTL for cached small files
}

func NewService(
	cache *cache.RedisCache,
	repo metadata.Repository,
	router objectstore.ObjectRoute,
	registry *objectstore.ProviderRegistry,
) *Service {
	return &Service{
		cache:              cache,
		metaRepo:           repo,
		router:             router,
		providers:          registry,
		smallFileThreshold: 1 * 1024 * 1024,
		cacheTTL:           24 * time.Hour,
	}
}

type PutRequest struct {
	Env           string
	LogicalRegion string
	Bucket        string
	Key           string
	ContentType   string
	Size          int64
	Body          io.Reader
	StorageClass  string // HOT/COLD/ARCHIVE
}

type PutResponse struct {
	ETag    string                `json:"etag"`
	Backend metadata.StoreBackend `json:"backend"`
	Size    int64                 `json:"size"`
}

func (s *Service) Put(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	if req.StorageClass == "" {
		req.StorageClass = "HOT"
	}
	if req.Size > 0 && req.Size <= s.smallFileThreshold {
		return s.putSmall(ctx, req)
	}
	return s.putLarge(ctx, req)
}

func (s *Service) putSmall(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, req.Body)
	if err != nil {
		return nil, err
	}
	data := buf.Bytes()
	cacheKey := s.cacheKey(req.Env, req.Bucket, req.Key)

	// 1. write to cache
	if err := s.cache.SetObject(ctx, cacheKey, data, s.cacheTTL); err != nil {
		// TODO: log warning
	}

	// 2. route to provider
	route, err := s.router.ResolveRoute(objectstore.RouteKey{
		Env:           req.Env,
		LogicalRegion: req.LogicalRegion,
		Bucket:        req.Bucket,
		StorageClass:  req.StorageClass,
	})
	if err != nil {
		return nil, err
	}

	backend, ok := s.providers.Get(route.ProviderName)
	if !ok {
		return nil, fmt.Errorf("no backend for provider %s", route.ProviderName)
	}

	loc := objectstore.ObjectLocation{
		ProviderType:   route.ProviderType,
		ProviderRegion: route.ProviderRegion,
		ProviderBucket: route.ProviderBucket,
		PhysicalKey:    s.buildPhysicalKey(req),
	}

	etag, err := backend.PutObject(ctx, loc, bytes.NewReader(data), n, objectstore.PutOptions{
		ContentType:  req.ContentType,
		StorageClass: req.StorageClass,
	})
	if err != nil {
		return nil, err
	}

	rec := &metadata.ObjectRecord{
		Env:            req.Env,
		LogicalRegion:  req.LogicalRegion,
		Bucket:         req.Bucket,
		ObjectKey:      req.Key,
		SizeBytes:      n,
		ContentType:    req.ContentType,
		StorageClass:   req.StorageClass,
		StoreBackend:   metadata.StoreRedisObject,
		ProviderType:   string(route.ProviderType),
		ProviderRegion: route.ProviderRegion,
		ProviderBucket: route.ProviderBucket,
		PhysicalKey:    loc.PhysicalKey,
		ETag:           etag,
		Status:         "ACTIVE",
	}
	if err := s.metaRepo.PutObject(ctx, rec); err != nil {
		return nil, err
	}

	return &PutResponse{
		ETag:    etag,
		Backend: rec.StoreBackend,
		Size:    n,
	}, nil
}

func (s *Service) putLarge(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	// Stream to object storage without caching.
	route, err := s.router.ResolveRoute(objectstore.RouteKey{
		Env:           req.Env,
		LogicalRegion: req.LogicalRegion,
		Bucket:        req.Bucket,
		StorageClass:  req.StorageClass,
	})
	if err != nil {
		return nil, err
	}
	backend, ok := s.providers.Get(route.ProviderName)
	if !ok {
		return nil, fmt.Errorf("no backend for provider %s", route.ProviderName)
	}
	loc := objectstore.ObjectLocation{
		ProviderType:   route.ProviderType,
		ProviderRegion: route.ProviderRegion,
		ProviderBucket: route.ProviderBucket,
		PhysicalKey:    s.buildPhysicalKey(req),
	}
	etag, err := backend.PutObject(ctx, loc, req.Body, req.Size, objectstore.PutOptions{
		ContentType:  req.ContentType,
		StorageClass: req.StorageClass,
	})
	if err != nil {
		return nil, err
	}

	rec := &metadata.ObjectRecord{
		Env:            req.Env,
		LogicalRegion:  req.LogicalRegion,
		Bucket:         req.Bucket,
		ObjectKey:      req.Key,
		SizeBytes:      req.Size,
		ContentType:    req.ContentType,
		StorageClass:   req.StorageClass,
		StoreBackend:   metadata.StoreObjectOnly,
		ProviderType:   string(route.ProviderType),
		ProviderRegion: route.ProviderRegion,
		ProviderBucket: route.ProviderBucket,
		PhysicalKey:    loc.PhysicalKey,
		ETag:           etag,
		Status:         "ACTIVE",
	}
	if err := s.metaRepo.PutObject(ctx, rec); err != nil {
		return nil, err
	}

	return &PutResponse{
		ETag:    etag,
		Backend: rec.StoreBackend,
		Size:    req.Size,
	}, nil
}

type GetRequest struct {
	Env           string
	LogicalRegion string
	Bucket        string
	Key           string
}

type GetResponse struct {
	Size        int64
	ContentType string
	Body        io.ReadCloser
}

func (s *Service) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	cacheKey := s.cacheKey(req.Env, req.Bucket, req.Key)

	// 1. try cache
	if data, err := s.cache.GetObject(ctx, cacheKey); err == nil && len(data) > 0 {
		return &GetResponse{
			Size:        int64(len(data)),
			ContentType: "", // in future we can cache meta as well
			Body:        io.NopCloser(bytes.NewReader(data)),
		}, nil
	}

	// 2. lookup metadata
	rec, err := s.metaRepo.GetObject(ctx, req.Env, req.LogicalRegion, req.Bucket, req.Key)
	if err != nil {
		return nil, err
	}

	routeName := rec.ProviderType // using provider type/name; here we treat ProviderType as key
	backend, ok := s.providers.Get(routeName)
	if !ok {
		// fallback: try by provider type if name-based lookup failed
		backend, ok = s.providers.Get(rec.ProviderBucket)
		if !ok {
			return nil, fmt.Errorf("no backend for provider %s", routeName)
		}
	}

	loc := objectstore.ObjectLocation{
		ProviderType:   objectstore.ProviderType(rec.ProviderType),
		ProviderRegion: rec.ProviderRegion,
		ProviderBucket: rec.ProviderBucket,
		PhysicalKey:    rec.PhysicalKey,
	}

	body, size, contentType, err := backend.GetObject(ctx, loc)
	if err != nil {
		return nil, err
	}

	// 3. optionally refill cache if small
	if size > 0 && size <= s.smallFileThreshold {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, body); err != nil {
			body.Close()
			return nil, err
		}
		data := buf.Bytes()
		_ = s.cache.SetObject(ctx, cacheKey, data, s.cacheTTL)
		body.Close()
		body = io.NopCloser(bytes.NewReader(data))
	}

	return &GetResponse{
		Size:        size,
		ContentType: contentType,
		Body:        body,
	}, nil
}

func (s *Service) cacheKey(env, bucket, key string) string {
	return fmt.Sprintf("obj:%s:%s:%s", env, bucket, key)
}

func (s *Service) buildPhysicalKey(req *PutRequest) string {
	// Simple physical key: env/logicalRegion/bucket/key
	return fmt.Sprintf("%s/%s/%s/%s", req.Env, req.LogicalRegion, req.Bucket, req.Key)
}
