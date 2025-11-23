package metadata

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type SQLRepository struct {
	conn *pgx.Conn
}

func NewSQLRepository(conn *pgx.Conn) *SQLRepository {
	return &SQLRepository{conn: conn}
}

func (r *SQLRepository) GetObject(ctx context.Context, env, region, bucket, key string) (*ObjectRecord, error) {
	const q = `
SELECT env, logical_region, bucket, object_key,
       size_bytes, content_type, storage_class, store_backend,
       provider_type, provider_region, provider_bucket, physical_key,
       etag, version, status, created_at, updated_at
FROM objects
WHERE env = $1 AND logical_region = $2 AND bucket = $3 AND object_key = $4 AND status = 'ACTIVE'
`
	row := r.conn.QueryRow(ctx, q, env, region, bucket, key)
	var rec ObjectRecord
	var storeBackend string
	if err := row.Scan(
		&rec.Env, &rec.LogicalRegion, &rec.Bucket, &rec.ObjectKey,
		&rec.SizeBytes, &rec.ContentType, &rec.StorageClass, &storeBackend,
		&rec.ProviderType, &rec.ProviderRegion, &rec.ProviderBucket, &rec.PhysicalKey,
		&rec.ETag, &rec.Version, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	rec.StoreBackend = StoreBackend(storeBackend)
	return &rec, nil
}

func (r *SQLRepository) PutObject(ctx context.Context, rec *ObjectRecord) error {
	if rec == nil {
		return errors.New("nil record")
	}
	now := time.Now()
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = now
	}
	rec.UpdatedAt = now
	if rec.Version == 0 {
		rec.Version = 1
	}
	if rec.Status == "" {
		rec.Status = "ACTIVE"
	}
	const q = `
INSERT INTO objects (
    env, logical_region, bucket, object_key,
    size_bytes, content_type, storage_class, store_backend,
    provider_type, provider_region, provider_bucket, physical_key,
    etag, version, status, created_at, updated_at
) VALUES (
    $1,$2,$3,$4,
    $5,$6,$7,$8,
    $9,$10,$11,$12,
    $13,$14,$15,$16,$17
)
ON CONFLICT (env, logical_region, bucket, object_key, status)
WHERE status = 'ACTIVE'
DO UPDATE SET
    size_bytes = EXCLUDED.size_bytes,
    content_type = EXCLUDED.content_type,
    storage_class = EXCLUDED.storage_class,
    store_backend = EXCLUDED.store_backend,
    provider_type = EXCLUDED.provider_type,
    provider_region = EXCLUDED.provider_region,
    provider_bucket = EXCLUDED.provider_bucket,
    physical_key = EXCLUDED.physical_key,
    etag = EXCLUDED.etag,
    version = objects.version + 1,
    updated_at = EXCLUDED.updated_at
`
	_, err := r.conn.Exec(ctx, q,
		rec.Env, rec.LogicalRegion, rec.Bucket, rec.ObjectKey,
		rec.SizeBytes, rec.ContentType, rec.StorageClass, string(rec.StoreBackend),
		rec.ProviderType, rec.ProviderRegion, rec.ProviderBucket, rec.PhysicalKey,
		rec.ETag, rec.Version, rec.Status, rec.CreatedAt, rec.UpdatedAt,
	)
	return err
}

func (r *SQLRepository) MarkDeleted(ctx context.Context, env, region, bucket, key string) error {
	const q = `
UPDATE objects
SET status = 'DELETED', updated_at = now()
WHERE env = $1 AND logical_region = $2 AND bucket = $3 AND object_key = $4 AND status = 'ACTIVE'
`
	cmd, err := r.conn.Exec(ctx, q, env, region, bucket, key)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
