# Firestore Data Access Layer (firestoredal)

The Firestore go client is already pretty easy to use but for simple use-cases it can be even easier. `firestoredal` helper streamlines these simple usage patterns into one small library.

### Import

```go
import "github.com/mchmarny/firestoredal"
```


### Helper Instance

```go

	ctx := context.Background()
	db, err := fsme.NewDB(ctx, PROJECT_ID, REGION)
	if err != nil {
		log.Panicf("Error while configuring DB: %v", err)
	}
	defer db.Close()

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