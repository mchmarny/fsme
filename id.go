package firestoredal

import (
	"fmt"
	"hash/fnv"

	uuid "github.com/satori/go.uuid"
)

const (
	idPrefix = "tid"
)

// GetNewID parses Firestore valid IDs (can't start with digits)
func GetNewID() string {
	return fmt.Sprintf("%s-%s", idPrefix, uuid.NewV4().String())
}

// ToID converts passed value to a valid Firestire ID
func ToID(val string) string {
	h := fnv.New32a()
	h.Write([]byte(val))
	return fmt.Sprintf("%s%d", idPrefix, h.Sum32())
}

// IsFavlidID validates that passed value is a valid Firestore ID
func IsValidID(val string) bool {
	
	if val == "" {
		return false
	}
	
	r := val[0]
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	
}