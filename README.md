# fsme

Simple Firestore DB helper

## Import

```go
import "github.com/mchmarny/fsme"
```

## Simple Save, Get, Delete

```go

	ctx := context.Background()
	db, err := NewDB(ctx, PROJECT_ID, REGION)
	if err != nil {
		log.Panicf("Error while configuring DB: %v", err)
	}

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