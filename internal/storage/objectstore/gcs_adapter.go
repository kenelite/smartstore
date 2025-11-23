package objectstore

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"time"
)

// GCSAdapter uses the native GCS client.
// Credentials are loaded via default application credentials or env GOOGLE_APPLICATION_CREDENTIALS,
// or can be controlled outside this code.
type GCSAdapter struct {
	client *storage.Client
}

func NewGCSAdapter(ctx context.Context) (*GCSAdapter, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSAdapter{client: client}, nil
}

func (a *GCSAdapter) PutObject(ctx context.Context, loc ObjectLocation, r io.Reader, size int64, opts PutOptions) (string, error) {
	wc := a.client.Bucket(loc.ProviderBucket).Object(loc.PhysicalKey).NewWriter(ctx)
	wc.ContentType = opts.ContentType
	if _, err := io.Copy(wc, r); err != nil {
		_ = wc.Close()
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	// GCS uses generation/etag; we can fetch attrs for ETag
	attrs, err := a.client.Bucket(loc.ProviderBucket).Object(loc.PhysicalKey).Attrs(ctx)
	if err != nil {
		return "", err
	}
	return attrs.Etag, nil
}

func (a *GCSAdapter) GetObject(ctx context.Context, loc ObjectLocation) (io.ReadCloser, int64, string, error) {
	rc, err := a.client.Bucket(loc.ProviderBucket).Object(loc.PhysicalKey).NewReader(ctx)
	if err != nil {
		return nil, 0, "", err
	}
	// GCS Reader knows the size
	return rc, rc.Attrs.Size, rc.ContentType(), nil
}

func (a *GCSAdapter) DeleteObject(ctx context.Context, loc ObjectLocation) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return a.client.Bucket(loc.ProviderBucket).Object(loc.PhysicalKey).Delete(ctx)
}

func (a *GCSAdapter) String() string {
	return fmt.Sprintf("GCSAdapter{%p}", a)
}
