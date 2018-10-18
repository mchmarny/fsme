package fsme

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

// DB represents simple FireStore helper
type DB struct {
	client    *firestore.Client
	projectID string
	region    string
	ctx       context.Context
}

// GetDBProjectID returns the GCP Project ID with which this DB was configured
func (d *DB) GetDBProjectID() string {
	return d.projectID
}

// GetDBRegion returns the GCP region with which this DB was configured
func (d *DB) GetDBRegion() string {
	return d.region
}

// Save inserts or updates by ID
func (d *DB) Save(collection string, obj *FSObject) error {

	if obj == nil {
		return errors.New("Object required")
	}

	if obj.ID == "" {
		return errors.New("Object ID required")
	}

	if collection == "" {
		return errors.New("Collection required")
	}

	_, err := d.client.Collection(collection).Doc(obj.ID).Set(d.ctx, obj)

	return err

}

// Get returns Firestore document for a given ID
func (d *DB) Get(collection, id string) (obj *FSObject, err error) {

	if id == "" {
		return nil, errors.New("ID required")
	}

	if collection == "" {
		return nil, errors.New("Collection required")
	}

	doc, err := d.client.Collection(collection).Doc(id).Get(d.ctx)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, fmt.Errorf("No data for ID: %s", id)
	}

	c := &FSObject{}

	if e := doc.DataTo(&c); err != nil {
		return nil, fmt.Errorf("Error parsing data: %v", e)
	}

	return c, nil

}

// Delete deletes Firestore document for a given ID
func (d *DB) Delete(collection, id string) error {

	if id == "" {
		return errors.New("ID required")
	}

	if collection == "" {
		return errors.New("Collection required")
	}

	_, err := d.client.Collection(collection).Doc(id).Delete(d.ctx)
	return err

}

// NewDB configures new DB instance
func NewDB(projectID, region string) (db *DB, err error) {

	if projectID == "" {
		return nil, errors.New("projectID variable required")
	}

	if region == "" {
		region = "us-central1"
		log.Printf("Region variable not set, using default: %s", region)
	}

	ctx := context.Background()
	c, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Error while creating Firestore client: %v", err)
	}

	d := &DB{
		ctx:       ctx,
		projectID: projectID,
		region:    region,
		client:    c,
	}

	return d, nil

}
