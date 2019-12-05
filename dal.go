package firestoredal

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/iterator"
)

// Save inserts or updates by ID
func (d *Store) Save(ctx context.Context, collection string, id string, obj interface{}) error {

	if obj == nil {
		return errors.New("object required")
	}

	if id == "" {
		return errors.New("id required")
	}

	if collection == "" {
		return errors.New("collection required")
	}

	_, err := d.client.Collection(collection).Doc(id).Set(ctx, obj)

	return err

}

// GetByID returns stored object for given ID
func (d *Store) GetByID(ctx context.Context, collection, id string, in interface{}) error {

	if id == "" {
		return errors.New("id required")
	}

	if collection == "" {
		return errors.New("collection required")
	}

	doc, err := d.client.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		return err
	}

	if doc == nil {
		return fmt.Errorf("no data for ID: %s", id)
	}

	if e := doc.DataTo(in); err != nil {
		return fmt.Errorf("error parsing data: %v", e)
	}

	return nil

}

// DeleteByID deletes stored object for a given ID
func (d *Store) DeleteByID(ctx context.Context, collection, id string) error {

	if id == "" {
		return errors.New("id required")
	}

	if collection == "" {
		return errors.New("collection required")
	}

	_, err := d.client.Collection(collection).Doc(id).Delete(ctx)
	return err

}

// QueryHandler provides new instance of an item
type QueryHandler interface {
	MakeNewItem() interface{}
	GetCriteria() []*StoreCriterion
	AddItem(item interface{})
}

// GetByQuery allows for filtered query using QueryHandler
func (d *Store) GetByQuery(ctx context.Context, collection string, handler QueryHandler) error {

	if handler == nil {
		return fmt.Errorf("criteria required")
	}

	q := d.client.Collection(collection).Query
	for _, c := range handler.GetCriteria() {
		q = q.Where(c.Property, c.Operator, c.Value)
	}

	docs := q.Documents(ctx)
	defer docs.Stop()

	for {
		d, e := docs.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			return e
		}

		item := handler.MakeNewItem()
		if e := d.DataTo(&item); e != nil {
			return e
		}
		handler.AddItem(item)
	}

	return nil

}
