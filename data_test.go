package fsme

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
	testCollectionName  = "test"
	testProjectID       = "s9-demo"
	testRegion          = ""
	numberOfTestRecords = 10
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

func TestGetAll(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestData")
	}
	ctx := context.Background()
	db, err := NewDB(ctx, testProjectID, testRegion)
	if err != nil {
		t.Errorf("Error while configuring DB: %v", err)
	}
	defer db.Close()

	// load objects
	list := []*FSObject{}
	savedCounter := 0

	for len(list) < numberOfTestRecords {

		savedCounter++
		obj := NewFSObject(map[string]interface{}{
			"RecordIndex": fmt.Sprintf("John-%d", savedCounter),
		})

		err = db.Save(testCollectionName, obj)
		if err != nil {
			t.Errorf("Error on save: %v", err)
		}

		list = append(list, obj)

	}

	objCh := make(chan *FSObject)
	receivedCounter := 0
	go func() {
		err = db.GetAll(testCollectionName, objCh)
		if err != nil {
			t.Errorf("Error on get: %v", err)
		}
	}()

	for {
		select {
		case objD := <-objCh:
			receivedCounter++
			t.Logf("Record: %d - %v", receivedCounter, objD.Data["RecordIndex"])
			err = db.Delete(testCollectionName, objD.ID)
			if err != nil {
				t.Errorf("Error on delete: %v", err)
			}
			if receivedCounter == numberOfTestRecords {
				return
			}
		default:
			// nothing to do here
		}
	}

}

func TestGetWhere(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping TestData")
	}
	ctx := context.Background()
	db, err := NewDB(ctx, testProjectID, testRegion)
	if err != nil {
		t.Errorf("Error while configuring DB: %v", err)
	}
	defer db.Close()

	// load objects
	list := []*FSObject{}

	obj1 := NewFSObject(map[string]interface{}{"City": "Portland"})
	db.Save(testCollectionName, obj1)
	list = append(list, obj1)

	obj2 := NewFSObject(map[string]interface{}{"City": "Seattle"})
	db.Save(testCollectionName, obj2)
	list = append(list, obj2)

	obj3 := NewFSObject(map[string]interface{}{"City": "Portland"})
	db.Save(testCollectionName, obj3)
	list = append(list, obj3)

	c := &FSCriterion{
		Property: "data.City",
		Operator: "==",
		Value:    "Portland",
	}

	list, err = db.GetWhere(testCollectionName, c)
	if err != nil {
		t.Errorf("Error on GetWhere: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("Got invalid number of records. Expected 2, Got %d", len(list))
	}

	// Delete
	db.Delete(testCollectionName, obj1.ID)
	db.Delete(testCollectionName, obj2.ID)
	db.Delete(testCollectionName, obj3.ID)

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
