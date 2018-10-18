package fsme

import (
	"time"
)

// FSObject represents simple object
type FSObject struct {
	ID   string      `json:"id" firestore:"id"`
	On   time.Time   `json:"created_on" firestore:"created_on"`
	Data interface{} `json:"data" firestore:"data"`
}

// NewFSObject returns fully loaded Firestore object
func NewFSObject(data interface{}) *FSObject {
	return &FSObject{
		ID:   getNewID(),
		On:   time.Now().UTC(),
		Data: data,
	}
}

// FirestoreValue is the payload of a FirestoreEvent event
type FirestoreValue struct {
	Fields interface{} `json:"fields"`
}

// FirestoreEvent is the Firestore document payload
type FirestoreEvent struct {
	OldValue FirestoreValue `json:"oldValue"`
	Value    FirestoreValue `json:"value"`
}
