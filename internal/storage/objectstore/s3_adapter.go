package objectstore

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Adapter is a production-ready baseline implementation using minio-go,
// which can talk to AWS S3 and Cloudflare R2 (S3-compatible).
type S3Adapter struct {
	client *minio.Client
}

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func NewS3Adapter(cfg S3Config) (*S3Adapter, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &S3Adapter{client: client}, nil
}

func (a *S3Adapter) PutObject(ctx context.Context, loc ObjectLocation, r io.Reader, size int64, opts PutOptions) (string, error) {
	putOpts := minio.PutObjectOptions{
		ContentType: opts.ContentType,
	}
	info, err := a.client.PutObject(ctx, loc.ProviderBucket, loc.PhysicalKey, r, size, putOpts)
	if err != nil {
		return "", err
	}
	return info.ETag, nil
}

func (a *S3Adapter) GetObject(ctx context.Context, loc ObjectLocation) (io.ReadCloser, int64, string, error) {
	obj, err := a.client.GetObject(ctx, loc.ProviderBucket, loc.PhysicalKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, "", err
	}
	stat, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, 0, "", err
	}
	return obj, stat.Size, stat.ContentType, nil
}

func (a *S3Adapter) DeleteObject(ctx context.Context, loc ObjectLocation) error {
	return a.client.RemoveObject(ctx, loc.ProviderBucket, loc.PhysicalKey, minio.RemoveObjectOptions{})
}
