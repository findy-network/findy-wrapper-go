package addons

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmuCfg_NewImmuCfg(t *testing.T) {
	cfg := NewImmuCfg(immuLedgerName)
	assert.NotNil(t, cfg)
}

func TestImmuCfg_Connect(t *testing.T) {
	cfg := NewImmuCfg(immuLedgerName)
	assert.NotNil(t, cfg)

	c, _, err := cfg.Connect()
	assert.NoError(t, err)

	err = c.Logout(context.Background())
	assert.NoError(t, err)
}
