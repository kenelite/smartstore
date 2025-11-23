CREATE TABLE IF NOT EXISTS objects (
  id              BIGSERIAL PRIMARY KEY,
  env             VARCHAR(16) NOT NULL,
  logical_region  VARCHAR(32) NOT NULL,
  bucket          VARCHAR(64) NOT NULL,
  object_key      TEXT NOT NULL,
  size_bytes      BIGINT NOT NULL,
  content_type    VARCHAR(255),
  storage_class   VARCHAR(32) NOT NULL,
  store_backend   VARCHAR(32) NOT NULL,
  provider_type   VARCHAR(32) NOT NULL,
  provider_region VARCHAR(32) NOT NULL,
  provider_bucket VARCHAR(255) NOT NULL,
  physical_key    TEXT NOT NULL,
  etag            VARCHAR(128),
  version         BIGINT NOT NULL DEFAULT 1,
  status          VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  created_at      TIMESTAMP NOT NULL DEFAULT now(),
  updated_at      TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_objects_active
ON objects (env, logical_region, bucket, object_key, status)
WHERE status = 'ACTIVE';
