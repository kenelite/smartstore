package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	apihttp "github.com/kenelite/smartstore/internal/api/http"
	"github.com/kenelite/smartstore/internal/cache"
	"github.com/kenelite/smartstore/internal/config"
	"github.com/kenelite/smartstore/internal/metadata"
	"github.com/kenelite/smartstore/internal/storage/objectstore"
	"github.com/kenelite/smartstore/internal/storage/smart"
)

func NewHTTPServer(cfg *config.Config) *http.Server {
	// init redis
	redisOpts := &redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	}
	redisCache := cache.NewRedisCache(redisOpts)

	// init metadata repository
	var repo metadata.Repository
	if cfg.DB.DSN != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := pgx.Connect(ctx, cfg.DB.DSN)
		if err != nil {
			log.Printf("failed to connect DB, fallback to in-memory repo: %v", err)
			repo = metadata.NewInMemoryRepository()
		} else {
			log.Printf("connected to DB")
			repo = metadata.NewSQLRepository(conn)
		}
	} else {
		log.Printf("no DB configured, using in-memory metadata repo")
		repo = metadata.NewInMemoryRepository()
	}

	// init router & provider registry
	route := objectstore.NewStaticRouter(cfg.ObjectStorage)
	registry := objectstore.NewProviderRegistry()

	// register providers based on config
	ctx := context.Background()
	for _, p := range cfg.ObjectStorage.Providers {
		switch p.Type {
		case config.ProviderAWS_S3, config.ProviderCF_R2:
			s3Adapter, err := objectstore.NewS3Adapter(objectstore.S3Config{
				Endpoint:  p.Endpoint,
				AccessKey: p.AccessKey,
				SecretKey: p.SecretKey,
				UseSSL:    p.UseSSL,
			})
			if err != nil {
				log.Printf("failed to init S3 adapter for provider %s: %v", p.Name, err)
				continue
			}
			registry.Register(p.Name, s3Adapter)
		case config.ProviderGCP_GCS:
			gcsAdapter, err := objectstore.NewGCSAdapter(ctx)
			if err != nil {
				log.Printf("failed to init GCS adapter for provider %s: %v", p.Name, err)
				continue
			}
			registry.Register(p.Name, gcsAdapter)
		default:
			log.Printf("unknown provider type %s for provider %s", p.Type, p.Name)
		}
	}

	smartSvc := smart.NewService(redisCache, repo, route, registry)
	handler := apihttp.NewHandler(smartSvc)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return srv
}
