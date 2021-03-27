package immu

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

	mock, isMock := c.client.(*mockImmuClient)
	loginCount := 0
	if isMock {
		loginCount = mock.loginOkCount
	}
	delta = 1 * time.Millisecond
	time.Sleep(5 * time.Millisecond)

	_, val, err = c.Read("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", val)
	if isMock {
		loginCount2 := mock.loginOkCount
		assert.Equal(t, loginCount+1, loginCount2)
	}
	c.Stop()
	time.Sleep(3 * time.Millisecond)
	immuLedgerImpl.Close()
}

func TestNeedRefresh(t *testing.T) {
	c := immuLedger
	c.loginTS = time.Now()
	delta = 4 * time.Millisecond
	time.Sleep(1 * time.Millisecond)
	assert.False(t, c.needRefresh())

	c.loginTS = time.Now()
	delta = 1 * time.Millisecond
	time.Sleep(4 * time.Millisecond)
	assert.True(t, c.needRefresh())
}
