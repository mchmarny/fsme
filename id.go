package firestoredal

import (
	"fmt"
	"hash/fnv"

	uuid "github.com/satori/go.uuid"
)

const (
	idPrefix = "tid"
)

// getNewID parses Firestore valid IDs (can't start with digits)
func getNewID() string {
	return fmt.Sprintf("%s-%s", idPrefix, uuid.NewV4().String())
}

func toID(query string) string {
	h := fnv.New32a()
	h.Write([]byte(query))
	return fmt.Sprintf("%s%d", idPrefix, h.Sum32())
}
