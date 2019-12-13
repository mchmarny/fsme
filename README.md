# lighter

[Cloud Firestore](https://cloud.google.com/firestore/) is a serverless data store with ACID transactions, multi-region replication, and powerful query engine. Its perfect for cloud-native solutions as it can autoscale with your application. Firestore has a very capable [golang client](https://godoc.org/cloud.google.com/go/firestore) (`cloud.google.com/go/firestore`) for reading and writing data to its database.

I created `lighter` as a wrapper to simplify many of the common (and sometimes verbose) tasks when using Firestore in my golang applications:

* GCP Project ID required to create Firestore golang client is automatically derived from your GCP metadata server or any of the [common environment variables](./meta.go) in case of local execution
* Common operations like `Save`, `GetByID`, and `Delete` are reduced to single line commands
* Firestore ID requirements are encapsulated into simple `GetNewID()` and `ToID(val)` operations (see [best practices](https://cloud.google.com/firestore/docs/best-practices))

Still, in case of more complex queries you will want to use the features of native golang Firestore client and `lighter` can still help.

* Create Firestore client without specifying GCP project ID (`NewClient` method)
* Simpler definition of your `where` and `order by` criteria (`QueryCriteria` struct, `GetQueryByCriteria` method)
* Generic method to process Firestore query results (`ResultHandler` interface, `HandleResults` method)

I hope you find `lighter` helpful. Issues and PRs are welcomed.

## Installation

To install `lighter` you will need Go version 1.11+ installed and execute the `go get` command:

```shell
go get -u github.com/mchmarny/lighter
```

Then in your code you can import `lighter`:

```go
import "github.com/mchmarny/lighter"
```

## Quick Start

Here is a simple application example assuming the following code is in your `main.go` file:

```go
package main

import (
	"context"
	"time"

	"github.com/mchmarny/lighter"
)

type Product struct {
	ID     string    `firestore:"id"`
	SoldOn time.Time `firestore:"sold"`
	Name   string    `firestore:"name"`
	Cost   float64    `firestore:"cost"`
}

func main() {
	ctx := context.Background()

	// create lighter client
	store, err := lighter.NewStore(ctx)
	handleError(err)
	defer store.Close()

	p1 := &Product{
		ID:     "id-1234", // Firestore IDs must start with a letter, see IDs section below
		SoldOn: time.Now().UTC(),
		Name:   "Demo Product",
		Cost:   2.99,
	}

	// save above defined product  to the product collection in Firestore
	err = store.Save(ctx, "product", p1.ID, p1)
	handleError(err)

	p2 := &Product{}
	// get a new instance (p2) of the above saved product (p1)
	err = store.GetByID(ctx, "product", p1.ID, p2)
	handleError(err)

	// delete the saved product
	err = store.DeleteByID(ctx, "product", p2.ID)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
```

## Get results sorted by struct property

> Use the name and case of the property defined in the struct `firestore` attribute

To get all products sorted by descending cost you will need to create instance of the `lighter.QueryCriteria` object and then pass it to the `store.GetByQuery` method:

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

## Get results filtered by struct property

> Use the name and case of the property defined in the struct `firestore` attribute

To get all products where the cost is less than 10.0 you need to created instance of the `lighter.QueryCriteria` object and pass it to the `store.GetByQuery` method:

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

`lighter` defines `ResultHandler` interface as a generic way to processing Firestore results:

```go
type ResultHandler interface {
	// MakeNew makes new instance of the saved struct
	MakeNew() interface{}
	// Append adds newly loaded item to results on each result iteration
	Append(item interface{})
}
```

In your code you will then need to define an implementation of that interface (e.g. the `ProductResultHandler`):

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

Once you have the struct that implements `ResultHandler` you can use it in either `GetByQuery`:

```go
h := &ProductResultHandler{Product: make([]*Product, 0)}
err := store.GetByQuery(ctx, q, h)
```

Or to handle results of your own Firestore documents query where you pass resulting `firestore.DocumentIterator` along with your handler:

```go
docs, err := client.Collection("products").Where("cost", ">-", 10.0)
handleError(err)

h := &ProductResultHandler{Product: make([]*Product, 0)}
err := store.HandleResults(ctx, docs, h)
```

## IDs

Firestore IDs must start with a letter. `lighter` provides a couple helpers in this area. You can either create brand new ID using the v4 UUID provider like this:

```go
id := lighter.GetNewID()
// results in something like this
// tid-7202525c-aa25-452c-a5c8-4d93fdc4074b
```

Or use existing value to create a valid ID. This is good approach to avoiding "hotspots" in your DB as these will eventually impact query latency.

```go
id := lighter.ToID("my-1234-value")
// results in 32-bit hash with the `tid-` prefix
```


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.


