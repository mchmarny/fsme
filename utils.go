package fsme

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

const (
	idPrefix = "tid"
)

// getNewID parses Firestore valid IDs (can't start with digits)
func getNewID() string {
	return fmt.Sprintf("%s-%s", idPrefix, uuid.NewV4().String())
}
