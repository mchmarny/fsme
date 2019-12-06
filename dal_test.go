package firestoredal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestStoreObject represents simple object
type TestStoreObject struct {
	ID   string                 `json:"id" firestore:"id"`
	On   time.Time              `json:"created_on" firestore:"created_on"`
	Data map[string]interface{} `json:"data" firestore:"data"`
}

type TestStoreObjectHandler struct {
	Items    []*TestStoreObject
	Criteria []*StoreCriterion
}

func (t *TestStoreObjectHandler) NewItem() interface{} {
	return &TestStoreObject{}
}

func (t *TestStoreObjectHandler) GetCriteria() []*StoreCriterion {
	return t.Criteria
}

func (t *TestStoreObjectHandler) HandleItem(item interface{}) {
	t.Items = append(t.Items, item.(*TestStoreObject))
}

// NewStoreObject returns fully loaded Firestore object
func NewStoreObject(data map[string]interface{}) *TestStoreObject {
	return &TestStoreObject{
		ID:   getNewID(),
		On:   time.Now().UTC(),
		Data: data,
	}
}

func TestGet(t *testing.T) {

	colName := "test_get"
	ctx := context.Background()
	store, err := NewStore(ctx)
	assert.Nil(t, err)

	data := map[string]interface{}{
		"Name":    "John",
		"Age":     40,
		"IsAdmin": false,
		"When":    time.Now().UTC(),
	}

	obj := NewStoreObject(data)
	err = store.Save(ctx, colName, obj.ID, obj)
	assert.Nil(t, err)

	obj2 := &TestStoreObject{}
	err = store.GetByID(ctx, colName, obj.ID, obj2)
	assert.Nil(t, err)
	assert.Equal(t, obj.ID, obj2.ID)

	err = store.DeleteByID(ctx, colName, obj.ID)
	assert.Nil(t, err)

	err = store.Close()
	assert.Nil(t, err)
}

func TestGetWhere(t *testing.T) {

	colName := "test_getwhere"
	ctx := context.Background()
	store, err := NewStore(ctx)
	assert.Nil(t, err)
	defer store.Close()

	// load objects
	list := make([]interface{}, 3)

	obj1 := NewStoreObject(map[string]interface{}{"City": "Portland"})
	store.Save(ctx, colName, obj1.ID, obj1)
	list[0] = obj1

	obj2 := NewStoreObject(map[string]interface{}{"City": "Seattle"})
	store.Save(ctx, colName, obj2.ID, obj2)
	list[1] = obj2

	obj3 := NewStoreObject(map[string]interface{}{"City": "Portland"})
	store.Save(ctx, colName, obj3.ID, obj3)
	list[2] = obj3

	h := &TestStoreObjectHandler{
		Criteria: []*StoreCriterion{
			&StoreCriterion{
				Property: "data.City",
				Operator: "==",
				Value:    "Portland",
			},
		},
		Items: make([]*TestStoreObject, 0),
	}

	err = store.GetByQuery(ctx, colName, h)
	assert.Nil(t, err)
	assert.NotNil(t, h)
	assert.Len(t, h.Items, 2)

	for i, o := range h.Items {
		t.Logf("obj[%d] = %v", i, o)
	}

	// Delete
	store.DeleteByID(ctx, colName, obj1.ID)
	store.DeleteByID(ctx, colName, obj2.ID)
	store.DeleteByID(ctx, colName, obj3.ID)

}

func TestNulData(t *testing.T) {

	colName := "test_getnil"
	ctx := context.Background()
	store, err := NewStore(ctx)
	assert.Nil(t, err)
	defer store.Close()

	obj := &TestStoreObject{}
	err = store.GetByID(ctx, colName, "invalidObjectID", obj)
	assert.NotNil(t, err)

	err = store.Close()
	assert.Nil(t, err)

}
