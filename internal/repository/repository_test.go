package repository

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRepository(t *testing.T) {
	t.Run("NewRepository", func(t *testing.T) {
		repo := NewRepository()
		require.NotNil(t, repo)
	})
}
