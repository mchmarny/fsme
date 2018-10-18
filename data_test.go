package fsme

import (
	"testing"
)

const (
	testCollectionName = "test"
)

type Person struct {
	Name string
	Age  int
}

func TestJobData(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestJobData")
	}

	db, err := NewDB("s9-demo", "")
	if err != nil {
		t.Errorf("Error while configuring DB: %v", err)
	}

	data := &Person{Name: "John", Age: 40}

	obj := NewFSObject(data)

	err = db.Save(testCollectionName, obj)

	if err != nil {
		t.Errorf("Error on save: %v", err)
	}

	obj2, err := db.Get(testCollectionName, obj.ID)

	if err != nil {
		t.Errorf("Error on get: %v", err)
	}

	if obj.ID != obj2.ID {
		t.Errorf("Got invalid data. Expected ID %s, Got ID: %s",
			obj.ID, obj2.ID)
	}

	err = db.Delete(testCollectionName, obj.ID)
	if err != nil {
		t.Errorf("Error on delete: %v", err)
	}

}
