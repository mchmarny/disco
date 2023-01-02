package object

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
)

// Save persists the provided content into bucket using name.
func Save(ctx context.Context, bucket, name, path string) error {
	if bucket == "" || name == "" || path == "" {
		return errors.New("bucket, name, and path are required")
	}

	d, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "error reading content from: %s", path)
	}
	if err := Put(ctx, bucket, name, d); err != nil {
		return errors.Wrapf(err, "error writing content from: %s to:%s/%s",
			path, bucket, name)
	}
	log.Printf("gs://%s/%v", bucket, name)

	return nil
}
