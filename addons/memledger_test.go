package addons

import (
	"testing"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/stretchr/testify/assert"
)

func TestMemLedger_Open(t *testing.T) {
	ok := memLedger.Open("FINDY_MEM_LEDGER")
	assert.True(t, ok)
}

func TestMemLedger_Write(t *testing.T) {
	ok := memLedger.Open("FINDY_MEM_LEDGER")
	assert.True(t, ok)
	err := memLedger.Write(plugin.TxDID, "testID", "testData")
	assert.NoError(t, err)
	name, value, err := memLedger.Read(plugin.TxDID, "testID")
	assert.NoError(t, err)
	assert.Equal(t, "testID", name)
	assert.Equal(t, "testData", value)
	err = memLedger.Write(plugin.TxDID, "testID2", "testData2")
	assert.NoError(t, err)
	name, value, err = memLedger.Read(plugin.TxDID, "testID2")
	assert.NoError(t, err)
	assert.Equal(t, "testID2", name)
	assert.Equal(t, "testData2", value)
}

func TestMemLedger_Read(t *testing.T) {
	ok := memLedger.Open("FINDY_MEM_LEDGER")
	assert.True(t, ok)
	err := memLedger.Write(plugin.TxDID, "testID", "testData")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		name, value, err := memLedger.Read(plugin.TxDID, "testID")
		assert.NoError(t, err)
		assert.Equal(t, "testID", name)
		assert.Equal(t, "testData", value)
	}
}
