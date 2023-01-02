package object

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	err := Save(context.Background(), "test", "", "")
	assert.Error(t, err)
	err = Save(context.Background(), "test", "test", "")
	assert.Error(t, err)
	err = Save(context.Background(), "test", "", "test")
	assert.Error(t, err)
	err = Save(context.Background(), "", "test", "test")
	assert.Error(t, err)
}
