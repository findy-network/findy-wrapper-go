package immu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestImmuCfg_NewImmuCfg(t *testing.T) {
	cfg := NewImmuCfg(immuLedgerName)
	assert.NotNil(t, cfg)
}

func TestImmuCfg_Connect(t *testing.T) {
	cfg := NewImmuCfg(immuLedgerName)
	assert.NotNil(t, cfg)

	c, token, err := cfg.Connect()
	assert.NoError(t, err)

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	err = c.Logout(ctx)
	assert.NoError(t, err)
}
