package firestoredal

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/iterator"
)

// DeleteByID deletes stored object for a given ID
func (d *Store) DeleteByID(ctx context.Context, collection, id string) error {

	if !IsValidID(id) {
		return fmt.Errorf("id must start with letter: '%s'", id)
	}

	if collection == "" {
		return errors.New("collection required")
	}

	_, err := d.client.Collection(collection).Doc(id).Delete(ctx)
	return err

}

// DeleteAll deletes all items in a collection
func (d *Store) DeleteAll(ctx context.Context, collection string, batchSize int) error {

	if collection == "" {
		return errors.New("collection required")
	}

	ref := d.client.Collection(collection)

	for {
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0
		batch := d.client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}

}
