package lighter

import (
	"context"
	"os"
	"testing"
)

var store *Store

func TestMain(m *testing.M) {
	ctx := context.Background()
	s, err := NewStore(ctx)
	if err != nil {
		panic(err)
	}
	store = s
	defer store.Close()
	os.Exit(m.Run())
}
