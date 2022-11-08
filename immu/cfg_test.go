package immu

import (
	"context"
	"testing"

	"github.com/lainio/err2/assert"
	"google.golang.org/grpc/metadata"
)

func TestImmuCfg_NewImmuCfg(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	cfg := NewImmuCfg(immuLedgerName)
	assert.INotNil(cfg)
}

func TestImmuCfg_Connect(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	cfg := NewImmuCfg(immuLedgerName)
	assert.INotNil(cfg)

	c, token, err := cfg.Connect()
	assert.NoError(err)

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	err = c.Logout(ctx)
	assert.NoError(err)
}
