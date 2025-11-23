package objectstore

import (
	"fmt"

	"github.com/kenelite/smartstore/internal/config"
)

type RouteKey struct {
	Env           string
	LogicalRegion string
	Bucket        string
	StorageClass  string
}

type RouteResult struct {
	ProviderName   string
	ProviderType   ProviderType
	ProviderRegion string
	ProviderBucket string
}

// ObjectRoute maps logical info to a physical provider/bucket.
type ObjectRoute interface {
	ResolveRoute(key RouteKey) (RouteResult, error)
}

type StaticRouter struct {
	routes []RouteResultWithKey
}

type RouteResultWithKey struct {
	Key    RouteKey
	Result RouteResult
}

func NewStaticRouter(cfg config.ObjectStorageConfig) *StaticRouter {
	rs := make([]RouteResultWithKey, 0, len(cfg.Routes))
	// Build provider lookup map by name
	providers := map[string]config.ProviderConfig{}
	for _, p := range cfg.Providers {
		providers[p.Name] = p
	}
	for _, r := range cfg.Routes {
		p, ok := providers[r.ProviderName]
		if !ok {
			continue
		}
		rs = append(rs, RouteResultWithKey{
			Key: RouteKey{
				Env:           r.Env,
				LogicalRegion: r.LogicalRegion,
				Bucket:        r.Bucket,
				StorageClass:  r.StorageClass,
			},
			Result: RouteResult{
				ProviderName:   p.Name,
				ProviderType:   ProviderType(p.Type),
				ProviderRegion: p.Region,
				ProviderBucket: r.ProviderBucket,
			},
		})
	}
	return &StaticRouter{routes: rs}
}

func (s *StaticRouter) ResolveRoute(key RouteKey) (RouteResult, error) {
	for _, r := range s.routes {
		if r.Key == key {
			return r.Result, nil
		}
	}
	return RouteResult{}, fmt.Errorf("no route for %+v", key)
}
