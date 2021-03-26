package addons

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	ok := immuLedger.Open("FINDY_IMMUDB_LEDGER")
	assert.True(t, ok)
	c := newClient(immuLedger)
	c.Start()

	err := c.Set("key1", "value1")
	assert.NoError(t, err)

	err = c.Set("key2", "value2")
	assert.NoError(t, err)

	c.realClient.ResetMemCache()

	val := ""
	val, err = c.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	val, err = c.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", val)

	c.Stop()
	time.Sleep(3 * time.Millisecond)
}
