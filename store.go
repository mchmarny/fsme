package lighter

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// Store represents simple FireStore helper
type Store struct {
	client *firestore.Client
}

// NewClient creates new Firestore client with derived project ID
func NewClient(ctx context.Context) (client *firestore.Client, err error) {

	if ctx == nil {
		return nil, errors.New("ctx required")
	}

	projectID, err := getProjectID()
	if err != nil {
		return nil, err
	}

	return firestore.NewClient(ctx, projectID)

}

// NewStore configures new client instance
func NewStore(ctx context.Context) (db *Store, err error) {

	c, err := NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	return &Store{
		client: c,
	}, nil

}

// NewStoreWithCredentialsFile configures new client instance with credentials file
func NewStoreWithCredentialsFile(ctx context.Context, path string) (db *Store, err error) {

	if ctx == nil {
		return nil, errors.New("ctx required")
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		return nil, fmt.Errorf("credential file does not exist: %s", path)
	}

	projectID, err := getProjectID()
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsFile(path)
	c, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		return nil, fmt.Errorf("error creating client with %s credential file: %v", path, err)
	}

	return &Store{
		client: c,
	}, nil

}

// Close closes client connection
func (d *Store) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}
