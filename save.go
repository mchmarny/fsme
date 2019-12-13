package firestoredal

import (
	"context"
	"errors"
	"fmt"
)

// Save inserts or updates by ID
func (d *Store) Save(ctx context.Context, collection string, id string, obj interface{}) error {

	if obj == nil {
		return errors.New("object required")
	}

	if !IsValidID(id) {
		return fmt.Errorf("id must start with letter: '%s'", id)
	}

	if collection == "" {
		return errors.New("collection required")
	}

	_, err := d.client.Collection(collection).Doc(id).Set(ctx, obj)

	return err

}
