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

import (
	"context"
	"time"

	"github.com/mchmarny/lighter"
)

type Product struct {
	ID     string    `json:"id" firestore:"id"`
	SoldOn time.Time `json:"sold" firestore:"sold"`
	Name   string    `json:"name" firestore:"name"`
	Cost   float64   `json:"cost" firestore:"cost"`
}

func main() {
	ctx := context.Background()
	store, err := lighter.NewStore(ctx)
	handleError(err)
	defer store.Close()

	p := &Product{
		ID:     "id-1234",
		SoldOn: time.Now().UTC(),
		Name:   "Demo Product",
		Cost:   2.99,
	}

	err = store.Save(ctx, "product", p.ID, p)
	handleError(err)

	p2 := &Product{}
	err = store.GetByID(ctx, "product", p.ID, p2)
	handleError(err)

	err = store.DeleteByID(ctx, "product", p2.ID)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
```

## Get sorted by struct property

> Use the name and case of the property defined in the struct `firestore` attribute

To get all products sorted by descending cost

```go
q := &lighter.QueryCriteria{
	Collection: "product",
	OrderBy: &lighter.Order{
		Property:   "cost",
		Descending: true,
	},
}

h := &ProductResultHandler{}
err = store.GetByQuery(ctx, q, h)
```

## Get filtered by struct property

> Use the name and case of the property defined in the struct `firestore` attribute

To get all products where the cost is less than 10.0

```go
q := &lighter.QueryCriteria{
	Collection: "product",
	Criteria: []*lighter.Criterion{
		&lighter.Criterion{
			Property: "cost",
			Operator: ">=",
			Value:    10.0,
		},
	},
}

h := &ProductResultHandler{}
err = store.GetByQuery(ctx, q, h)
```

## Process results using custom handler

`lighter` defines `ResultHandler` interface to process Firestore results

```go
// ResultHandler defines methods required to handle result items
type ResultHandler interface {
	// MakeNew makes new item instance for loading from result iterator
	MakeNew() interface{}
	// Append adds newly loaded item to results
	Append(item interface{})
}
```

In your code you need to define an implementation of that interface (e.g. the `ProductResultHandler`)

```go
type ProductResultHandler struct {
	Products []*Product
}

func (t *ProductResultHandler) MakeNew() interface{} {
	return &Product{}
}

func (t *ProductResultHandler) Append(item interface{}) {
	t.Products = append(t.Products, item.(*Product))
}
```

Once you have the struct that implements `ResultHandler` you can use it either in `GetByQuery`

```go
h := &ProductResultHandler{Product: make([]*Product, 0)}
err := store.GetByQuery(ctx, q, h)
```

Or by creating your Firestore documents query and passing the resulting `firestore.DocumentIterator` along with your handler

```go
docs, err := client.Collection("products").Where("cost", ">-", 10.0)
handleError(err)

h := &ProductResultHandler{Product: make([]*Product, 0)}
err := store.HandleResults(ctx, docs, h)
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


