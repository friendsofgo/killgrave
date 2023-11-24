package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImposterFS(t *testing.T) {
	t.Run("imposters directory not found", func(t *testing.T) {
		_, err := NewImposterFS("failImposterPath")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "the directory failImposterPath doesn't exists")
	})

	t.Run("existing imposters directory", func(t *testing.T) {
		_, err := NewImposterFS("test/testdata/imposters")
		assert.NoError(t, err)
	})
}
