package lighter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	store, err := NewStore(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, store)
}

func TestNewStoreWithCredentialsAndNoFile(t *testing.T) {
	ctx := context.Background()
	store, err := NewStoreWithCredentialsFile(ctx, "no-file")
	assert.NotNil(t, err)
	assert.Nil(t, store)
}
