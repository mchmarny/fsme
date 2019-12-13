# lighter

`lighter` is the small [Cloud Firestore](https://cloud.google.com/firestore/) client wrapper to simplify many of the common (and sometimes verbose) tasks when using Firestore in your go application. `lighter` can also be used with the official go client library to provide simpler `Query` definition or `ResultHandler` to help you extract Firestore documents to go structures.

## Installation

To install `lighter` you will need Go version 1.11+ installed

```shell
go get -u github.com/mchmarny/lighter
```

Then in your code you can import `lighter` like this

```go
import "github.com/mchmarny/lighter"
```

## Quick Start

Assuming the following code is in your `main.go` file

```go
package main

import "github.com/mchmarny/lighter"

type Product struct {
	ID    	string    `json:"id" firestore:"id"`
	SoldOn  time.Time `json:"sold" firestore:"sold"`
	Name    string    `json:"name" firestore:"name"`
	Cost    float64    `json:"cost" firestore:"cost"`
}

func main() {
	ctx := context.Background()
	store, err := lighter.NewStore(ctx)
	handleError(err)
	defer store.Close()

	p := &Product{
		ID:    "id-1234",
		SoldOn:    time.Now().UTC(),
		Name:  "Demo Product",
		Cost: 2.99,
	}

	err = store.Save(ctx, "product", p.ID, p)
	handleError(err)

	p2 := &Product{}
	err = store.GetByID(ctx, "product", p.ID, p2)
	handleError(err)

	err = store.DeleteByID(ctx, colName, obj.ID)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
```

> Define `FS_CLIENT_IDENTITY` environment variable with a path to your GCP service account file to have the client use specific identity

### Simple Save, Get, Delete

```go

	// Your data
	myData := map[string]interface{}{
		"Name":    "John",
		"Age":     40,
		"IsAdmin": false,
		"When":    time.Now().UTC(),
	}

	// DB record
	dbObj := fsme.NewFSObject(myData)

	// Save
	err = db.Save("users", dbObj)
	if err != nil {
		log.Panicf("Error on save: %v", err)
	}

	// Get
	obj2, err := db.Get("users", dbObj.ID)
	if err != nil {
		log.Panicf("Error on get: %v", err)
	}

	// Delete
	err = db.Delete("users", obj2.ID)
	if err != nil {
		log.Panicf("Error on delete: %v", err)
	}

	// Finish
	err = db.Close()
	if err != nil {
		log.Panicf("Error on close: %v", err)
	}

```

### Get All


```go

	objCh := make(chan *fsme.FSObject)
	go func() {
		err = db.GetAll("users", objCh)
		if err != nil {
			log.Panicf("Error on get: %v", err)
		}
	}()

	for {
		select {
		case rec := <-objCh:
			log.Printf("Record: %v", rec.ID)
		default:
			// nothing to do here
		}
	}

```

### Get All Where

```go

	c := &fsme.FSCriterion{
		Property: "data.City",
		Operator: "==",
		Value:    "Portland",
	}

	list, err = db.GetWhere("users", c)
	if err != nil {
		log.Panicf("Error on get where: %v", err)
	}

	for i, u := range list {
		log.Printf("User[%d] - %v", i, u)
	}


```