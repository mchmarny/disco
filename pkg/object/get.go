package object

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

const (
	timeoutSeconds = 60
)

var (
	// ErrObjectNotFound is returned when object does not exists.
	ErrObjectNotFound = errors.New("object not found")
)

// Get retrieves content from bucket using name or returns the
// ObjectNotFound if does not exists.
func Get(ctx context.Context, bucket, name string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating storage client")
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*timeoutSeconds)
	defer cancel()

	b := client.Bucket(bucket)

	it := b.Objects(ctx, &storage.Query{Prefix: name})
	_, err = it.Next()
	if errors.Is(err, iterator.Done) {
		return nil, ErrObjectNotFound
	}

	rc, err := b.Object(name).NewReader(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error creating reader")
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "error reading data")
	}

	return data, nil
}
