package immu

import (
	"testing"
	"time"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/lainio/err2/assert"
)

func TestSetAndGet(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	ok := immuLedgerImpl.Open("FINDY_IMMUDB_LEDGER")
	assert.That(ok)
	c := newClient(immuLedgerImpl)
	c.Start()

	err := c.Write(plugin.TxDID, "key1", "value1")
	assert.NoError(err)

	err = c.Write(plugin.TxDID, "key2", "value2")
	assert.NoError(err)

	c.ResetMemCache()

	val := ""
	_, val, err = c.Read(plugin.TxDID, "key1")
	assert.NoError(err)
	assert.Equal("value1", val)

	mock, isMock := c.client.(*mockImmuClient)
	loginCount := 0
	if isMock {
		loginCount = mock.loginOkCount
	}
	delta = 1 * time.Millisecond
	time.Sleep(5 * time.Millisecond)

	_, val, err = c.Read(plugin.TxDID, "key2")
	assert.NoError(err)
	assert.Equal("value2", val)
	if isMock {
		loginCount2 := mock.loginOkCount
		assert.Equal(loginCount+1, loginCount2)
	}
	c.Stop()
	time.Sleep(3 * time.Millisecond)
	immuLedgerImpl.Close()
}

func TestNeedRefresh(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	c := immuLedger
	c.loginTS = time.Now()
	delta = 4 * time.Millisecond
	time.Sleep(1 * time.Millisecond)
	assert.ThatNot(c.needRefresh())

	c.loginTS = time.Now()
	delta = 1 * time.Millisecond
	time.Sleep(4 * time.Millisecond)
	assert.That(c.needRefresh())
}
