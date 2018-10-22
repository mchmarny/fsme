package fsme

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// DB represents simple FireStore helper
type DB struct {
	client *firestore.Client
	config GCPConfig
	ctx    context.Context
}

// GetConfig returns the GCP info with which this DB was configured
func (d *DB) GetConfig() GCPConfig {
	return d.config
}

// Close closes client connection
func (d *DB) Close() error {
	return d.client.Close()
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

// GetAll gets all records in a specified collection
func (d *DB) GetAll(collection string, ch chan<- *FSObject) error {

	i := d.client.Collection(collection).Documents(d.ctx)
	defer i.Stop()

	for {
		d, e := i.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			return e
		}

		c := &FSObject{}
		if e := d.DataTo(&c); e != nil {
			return e
		}

		if ch != nil {
			ch <- c
		}

	}

	return nil

}

// GetWhere allows for filtered query using FSWhereCriterion object
func (d *DB) GetWhere(collection string, c *FSCriterion) (list []*FSObject, err error) {

	if c == nil {
		return nil, fmt.Errorf("Criteria required")
	}

	result := make([]*FSObject, 0)
	i := d.client.Collection(collection).
		Where(c.Property, c.Operator, c.Value).
		Documents(d.ctx)

	defer i.Stop()

	for {
		d, e := i.Next()
		if e == iterator.Done {
			break
		}
		if e != nil {
			return result, e
		}

		c := &FSObject{}
		if e := d.DataTo(&c); e != nil {
			return result, e
		}

		result = append(result, c)

	}

	return result, nil

}

// NewDB configures new DB instance
func NewDB(ctx context.Context, projectID, region string) (db *DB, err error) {

	if ctx == nil {
		return nil, errors.New("ctx required")
	}

	if projectID == "" {
		return nil, errors.New("projectID required")
	}

	if region == "" {
		region = "us-central1"
		log.Printf("Region not set, using default: %s", region)
	}

	conf := GCPConfig{
		ProjectID: projectID,
		Region:    region,
		SetOn:     time.Now().UTC(),
	}

	// Connection options
	var c *firestore.Client

	clientIdentity := os.Getenv("FS_CLIENT_IDENTITY")
	if clientIdentity != "" {
		log.Printf("Using credentials file: %s", clientIdentity)
		opt := option.WithCredentialsFile(clientIdentity)
		c, err = firestore.NewClient(ctx, conf.ProjectID, opt)
	} else {
		log.Print("No credentials defined, using defaults")
		c, err = firestore.NewClient(ctx, conf.ProjectID)
	}

	if err != nil {
		return nil, fmt.Errorf("Error while creating Firestore client: %v", err)
	}

	d := &DB{
		ctx:    ctx,
		config: conf,
		client: c,
	}

	log.Printf("DB configured %s", d.GetConfig())

	return d, nil

}
