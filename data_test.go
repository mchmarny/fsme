package fsme

import (
	"testing"
	"time"
)

const (
	testCollectionName = "test"
)

func TestJobData(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestJobData")
	}

	db, err := NewDB("s9-demo", "")
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

}
