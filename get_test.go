package lighter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {

	colName := "test_get"
	ctx := context.Background()
	store, err := NewStore(ctx)
	assert.Nil(t, err)
	defer store.Close()

	obj := NewTestObject("John", 40, 2.75)
	err = store.Save(ctx, colName, obj.ID, obj)
	assert.Nil(t, err)

	obj2 := &MockedStoreObject{}
	err = store.GetByID(ctx, colName, obj.ID, obj2)
	assert.Nil(t, err)
	assert.Equal(t, obj.ID, obj2.ID)

	err = store.DeleteByID(ctx, colName, obj.ID)
	assert.Nil(t, err)

}

func TestQuerySort(t *testing.T) {

	colName := "test_sortcol"
	ctx := context.Background()
	store, err := NewStore(ctx)
	assert.Nil(t, err)
	err = store.DeleteAll(ctx, colName, 1)
	assert.Nil(t, err)
	defer store.Close()

	obj1 := NewTestObject("B", 1, 0.1)
	store.Save(ctx, colName, obj1.ID, obj1)

	obj2 := NewTestObject("B", 2, 0.2)
	store.Save(ctx, colName, obj2.ID, obj2)

	obj3 := NewTestObject("B", 3, 0.3)
	store.Save(ctx, colName, obj3.ID, obj3)

	obj4 := NewTestObject("C", 4, 0.4)
	store.Save(ctx, colName, obj4.ID, obj4)

	obj5 := NewTestObject("B", 5, 0.5)
	store.Save(ctx, colName, obj5.ID, obj5)

	q := &QueryCriteria{
		Collection: colName,
		OrderBy: &Order{
			Property:   "name",
			Descending: true,
		},
	}

	h := &TestObjectHandler{
		Items: make([]*MockedStoreObject, 0),
	}

	err = store.GetByQuery(ctx, q, h)
	assert.Nil(t, err)

	// item with the lowest count is first
	assert.NotNil(t, h.Items)
	assert.NotNil(t, h.Items[0])
	assert.Equal(t, "C", h.Items[0].Name)

}

func TestGetByQuery(t *testing.T) {

	colName := "test_getcriterion"
	ctx := context.Background()
	store, err := NewStore(ctx)
	assert.Nil(t, err)
	err = store.DeleteAll(ctx, colName, 1)
	assert.Nil(t, err)
	defer store.Close()

	obj1 := NewTestObject("Portland", 1, 0.1)
	store.Save(ctx, colName, obj1.ID, obj1)

	obj2 := NewTestObject("Seattle", 2, 0.2)
	store.Save(ctx, colName, obj2.ID, obj2)

	obj3 := NewTestObject("Portland", 3, 0.3)
	store.Save(ctx, colName, obj3.ID, obj3)

	q := &QueryCriteria{
		Collection: colName,
		Criteria: []*Criterion{
			&Criterion{
				Property: "name",
				Operator: "==",
				Value:    "Seattle",
			},
		},
	}

	h := &TestObjectHandler{
		Items: make([]*MockedStoreObject, 0),
	}

	err = store.GetByQuery(ctx, q, h)
	assert.Nil(t, err)

	// handler has expected number of items (2 Portlands)
	assert.NotNil(t, h)
	assert.Len(t, h.Items, 1)

	// item with the lowest count is first
	assert.NotNil(t, h.Items)
	assert.NotNil(t, h.Items[0])
	assert.Equal(t, 2, h.Items[0].Count)

}

func TestNulData(t *testing.T) {

	colName := "test_getnil"
	ctx := context.Background()
	store, err := NewStore(ctx)
	assert.Nil(t, err)
	defer store.Close()

	obj := &MockedStoreObject{}
	err = store.GetByID(ctx, colName, "invalidObjectID", obj)
	assert.NotNil(t, err)

	err = store.Close()
	assert.Nil(t, err)

}
