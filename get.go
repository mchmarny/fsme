package lighter

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// GetByID returns stored object for given ID
func (d *Store) GetByID(ctx context.Context, collection, id string, in interface{}) error {

	if !IsValidID(id) {
		return fmt.Errorf("id must start with letter: '%s'", id)
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

// ResultHandler defines methods required to handle result items
type ResultHandler interface {
	// MakeNew makes new item instance for loading from result iterator
	MakeNew() interface{}
	// Append adds newly loaded item to results
	Append(item interface{})
}

// QueryCriteria defines the Firestore query query
type QueryCriteria struct {
	Collection string
	Criteria   []*Criterion
	OrderBy    *Order
}

// Order defines a single Firestore property sort order
type Order struct {
	Property   string
	Descending bool
}

// Criterion defines single Firestore where criteria
type Criterion struct {
	// Property is the name of the property in where clause. Assumes index in Firestore
	Property string
	Operator string
	Value    interface{}
}

// GetByQuery allows for filtered query using QueryHandler
func (d *Store) GetByQuery(ctx context.Context, q *QueryCriteria, h ResultHandler) error {

	if q == nil {
		return fmt.Errorf("query required")
	}

	if h == nil {
		return fmt.Errorf("handler required")
	}

	sq, err := GetQueryByCriteria(d.client, q)
	if err != nil {
		return fmt.Errorf("error building query: %v", err)
	}

	docs := sq.Documents(ctx)
	defer docs.Stop()

	return HandleResults(ctx, docs, h)

}

// GetQueryByCriteria builds Firestore query
func GetQueryByCriteria(c *firestore.Client, q *QueryCriteria) (query *firestore.Query, err error) {

	if c == nil {
		return nil, fmt.Errorf("client required")
	}

	if q == nil {
		return nil, fmt.Errorf("query required")
	}

	sq := c.Collection(q.Collection).Query

	if q.Criteria != nil {
		for _, c := range q.Criteria {
			sq = sq.Where(c.Property, c.Operator, c.Value)
		}
	}

	if q.OrderBy != nil {
		dir := firestore.Asc
		if q.OrderBy.Descending {
			dir = firestore.Desc
		}
		sq = sq.OrderBy(q.OrderBy.Property, dir)
	}

	return &sq, nil

}

// HandleResults allows for filtered query using QueryHandler
func HandleResults(ctx context.Context, docs *firestore.DocumentIterator, h ResultHandler) error {

	if docs == nil {
		return fmt.Errorf("doc iterator required")
	}

	if h == nil {
		return fmt.Errorf("handler required")
	}

	for {
		d, e := docs.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			return e
		}

		item := h.MakeNew()
		if e := d.DataTo(&item); e != nil {
			return e
		}
		h.Append(item)
	}

	return nil

}
