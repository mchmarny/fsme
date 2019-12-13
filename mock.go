package lighter

import (
	"time"
)

// MockedStoreObject represents simple object
type MockedStoreObject struct {
	ID    string    `json:"id" firestore:"id"`
	On    time.Time `json:"on" firestore:"on"`
	Name  string    `json:"name" firestore:"name"`
	Count int       `json:"count" firestore:"count"`
	Value float64   `json:"value" firestore:"value"`
}

// NewTestObject returns fully loaded Firestore object
func NewTestObject(name string, count int, value float64) *MockedStoreObject {
	return &MockedStoreObject{
		ID:    GetNewID(),
		On:    time.Now().UTC(),
		Name:  name,
		Count: count,
		Value: value,
	}
}

// TestObjectHandler is a test implementation of the QueryHandler interface
type TestObjectHandler struct {
	Items []*MockedStoreObject
}

// MakeNew makes new item instance for loading from result iterator
func (t *TestObjectHandler) MakeNew() interface{} {
	return &MockedStoreObject{}
}

// Append adds newly loaded item to results
func (t *TestObjectHandler) Append(item interface{}) {
	t.Items = append(t.Items, item.(*MockedStoreObject))
}
