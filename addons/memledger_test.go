package addons

import (
	"testing"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/lainio/err2/assert"
)

func TestMemLedger_Open(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	ok := memLedger.Open("FINDY_MEM_LEDGER")
	assert.That(ok)
}

func TestMemLedger_Write(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	ok := memLedger.Open("FINDY_MEM_LEDGER")
	assert.That(ok)
	err := memLedger.Write(plugin.TxDID, "testID", "testData")
	assert.NoError(err)
	err = memLedger.Write(plugin.TxDID, "testID", "testData")
	assert.Error(err)
	name, value, err := memLedger.Read(plugin.TxDID, "testID")
	assert.NoError(err)
	assert.Equal("testID", name)
	assert.Equal("testData", value)
	err = memLedger.Write(plugin.TxDID, "testID2", "testData2")
	assert.NoError(err)
	name, value, err = memLedger.Read(plugin.TxDID, "testID2")
	assert.NoError(err)
	assert.Equal("testID2", name)
	assert.Equal("testData2", value)
}

func TestMemLedger_Read(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	ok := memLedger.Open("FINDY_MEM_LEDGER")
	assert.That(ok)
	err := memLedger.Write(plugin.TxDID, "testID", "testData")
	assert.NoError(err)

	for i := 0; i < 100; i++ {
		name, value, err := memLedger.Read(plugin.TxDID, "testID")
		assert.NoError(err)
		assert.Equal("testID", name)
		assert.Equal("testData", value)
	}
}
