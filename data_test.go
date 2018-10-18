package fsme

import (
	"context"
	"testing"
	"time"
)

const (
	testCollectionName = "test"
	testProjectID      = "s9-demo"
	testRegion         = ""
)

func TestData(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestData")
	}
	ctx := context.Background()
	db, err := NewDB(ctx, testProjectID, testRegion)
	if err != nil {
		t.Errorf("Error while configuring DB: %v", err)
	}

	data := map[string]interface{}{
		"Name":    "John",
		"Age":     40,
		"IsAdmin": false,
		"When":    time.Now().UTC(),
	}

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

	if obj.Data["Name"] != obj2.Data["Name"] {
		t.Errorf("Got invalid data. Expected Name %s, Got: %s",
			obj.Data["Name"], obj2.Data["Name"])
	}

	err = db.Delete(testCollectionName, obj.ID)
	if err != nil {
		t.Errorf("Error on delete: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Error on close: %v", err)
	}
}

func TestNulData(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestNulData")
	}
	ctx := context.Background()
	db, err := NewDB(ctx, testProjectID, testRegion)
	if err != nil {
		t.Errorf("Error while configuring DB: %v", err)
	}

	obj, err := db.Get(testCollectionName, "invalidObjectID")

	if err == nil {
		t.Error("Expected error")
	}

	if obj != nil {
		t.Error("Got data for an invalid object")
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Error on close: %v", err)
	}

}
