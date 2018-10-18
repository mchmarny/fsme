package fsme

import (
	"fmt"
	"time"
)

// FSObject represents simple object
type FSObject struct {
	ID   string                 `json:"id" firestore:"id"`
	On   time.Time              `json:"created_on" firestore:"created_on"`
	Data map[string]interface{} `json:"data" firestore:"data"`
}

// GCPConfig represents GCP config
type GCPConfig struct {
	ProjectID string
	Region    string
	SetOn     time.Time
}

func (c GCPConfig) String() string {
	return fmt.Sprintf("[ ProjectID:%s, Region:%s, SetOn:%v ]",
		c.ProjectID, c.Region, c.SetOn)
}

// NewFSObject returns fully loaded Firestore object
func NewFSObject(data map[string]interface{}) *FSObject {
	return &FSObject{
		ID:   getNewID(),
		On:   time.Now().UTC(),
		Data: data,
	}
}
