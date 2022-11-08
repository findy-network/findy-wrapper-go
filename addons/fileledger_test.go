package addons

import (
	"testing"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/lainio/err2/assert"
)

func TestFileLedger_Open(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	ok := fileLedger.Open("FINDY_FILE_LEDGER_TEST")
	assert.That(ok)
}

func TestFileLedger_Write(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	ok := fileLedger.Open("FINDY_FILE_LEDGER_TEST")
	assert.That(ok)
	err := fileLedger.Write(plugin.TxDID, "testID", "testData")
	assert.NoError(err)
	name, value, err := fileLedger.Read(plugin.TxDID, "testID")
	assert.NoError(err)
	assert.Equal("testID", name)
	assert.Equal("testData", value)
	err = fileLedger.Write(plugin.TxDID, "testID2", "testData2")
	assert.NoError(err)
	name, value, err = fileLedger.Read(plugin.TxDID, "testID2")
	assert.NoError(err)
	assert.Equal("testID2", name)
	assert.Equal("testData2", value)
}

func TestFileLedger_Read(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
	ok := fileLedger.Open("FINDY_FILE_LEDGER_TEST")
	assert.That(ok)
	err := fileLedger.Write(plugin.TxDID, "testID", "testData")
	assert.NoError(err)

	for i := 0; i < 100; i++ {
		name, value, err := fileLedger.Read(plugin.TxDID, "testID")
		assert.NoError(err)
		assert.Equal("testID", name)
		assert.Equal("testData", value)
	}
}
