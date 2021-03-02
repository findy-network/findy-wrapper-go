package addons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLedger_Open(t *testing.T) {
	ok := fileLedger.Open("FINDY_FILE_LEDGER_TEST")
	assert.True(t, ok)
}

func TestFileLedger_Write(t *testing.T) {
	ok := fileLedger.Open("FINDY_FILE_LEDGER_TEST")
	assert.True(t, ok)
	err := fileLedger.Write("testID", "testData")
	assert.NoError(t, err)
	name, value, err := fileLedger.Read("testID")
	assert.NoError(t, err)
	assert.Equal(t, "testID", name)
	assert.Equal(t, "testData", value)
	err = fileLedger.Write("testID2", "testData2")
	assert.NoError(t, err)
	name, value, err = fileLedger.Read("testID2")
	assert.NoError(t, err)
	assert.Equal(t, "testID2", name)
	assert.Equal(t, "testData2", value)
}

func TestFileLedger_Read(t *testing.T) {
	ok := fileLedger.Open("FINDY_FILE_LEDGER_TEST")
	assert.True(t, ok)
	err := fileLedger.Write("testID", "testData")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		name, value, err := fileLedger.Read("testID")
		assert.NoError(t, err)
		assert.Equal(t, "testID", name)
		assert.Equal(t, "testData", value)
	}
}