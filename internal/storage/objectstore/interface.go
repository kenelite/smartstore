package objectstore

import (
	"context"
	"io"
)

type ProviderType string

const (
	ProviderAWS_S3  ProviderType = "AWS_S3"
	ProviderCF_R2   ProviderType = "CF_R2"
	ProviderGCP_GCS ProviderType = "GCP_GCS"
)

type ObjectLocation struct {
	ProviderType   ProviderType
	ProviderRegion string
	ProviderBucket string
	PhysicalKey    string
}

type PutOptions struct {
	ContentType  string
	Metadata     map[string]string
	StorageClass string // HOT/COLD/ARCHIVE
}

type ObjectStorage interface {
	PutObject(ctx context.Context, loc ObjectLocation, r io.Reader, size int64, opts PutOptions) (etag string, err error)
	GetObject(ctx context.Context, loc ObjectLocation) (body io.ReadCloser, size int64, contentType string, err error)
	DeleteObject(ctx context.Context, loc ObjectLocation) error
}
