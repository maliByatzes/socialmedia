package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
  require.NoError(t, os.Setenv("CLIENT_URL", "http://localhost:3000"))
  require.NoError(t, os.Setenv("DB_URL", "database_url"))
  require.NoError(t, os.Setenv("PORT", "6969"))

  cfg, err := NewConfig()
  require.NoError(t, err)

  assert.Equal(t, cfg.ClientURL, "http://localhost:3000")
  assert.Equal(t, cfg.DBURL, "database_url")
  assert.Equal(t, cfg.Port, "6969")
}
