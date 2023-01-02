package object

import (
	"bytes"
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

// Put persists the provided content into bucket using name.
func Put(ctx context.Context, bucket, name string, content []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error creating storage client")
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*timeoutSeconds)
	defer cancel()

	wc := client.Bucket(bucket).Object(name).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, bytes.NewReader(content)); err != nil {
		return errors.Wrap(err, "error copping content to writer")
	}

	if err := wc.Close(); err != nil {
		return errors.Wrap(err, "error closing writer")
	}

	return nil
}
