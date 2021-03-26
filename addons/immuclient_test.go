package addons

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	ok := immuLedgerImpl.Open("FINDY_IMMUDB_LEDGER")
	assert.True(t, ok)
	c := newClient(immuLedgerImpl)
	c.Start()

	err := c.Write("key1", "value1")
	assert.NoError(t, err)

	err = c.Write("key2", "value2")
	assert.NoError(t, err)

	c.ResetMemCache()

	val := ""
	_, val, err = c.Read("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	_, val, err = c.Read("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", val)

	c.Stop()
	time.Sleep(3 * time.Millisecond)
}
