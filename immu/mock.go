package immu

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/codenotary/immudb/pkg/api/schema"
	im "github.com/codenotary/immudb/pkg/client"
	"github.com/golang/glog"
)

type keyType = string

var store = make(map[keyType][]byte)

// mockImmuClient is a mock for the im.ImmuClient interface. We implement only
// subset of the methods of the full interface. We MUST implement all of them
// we are calling from the addon! If not function ptr is nil and crash.
type mockImmuClient struct {
	im.ImmuClient // mocked interface (full version)

	setOkCount   int // incremented on every call of the Set()
	getOkCount   int // incremented on every call of the Get()
	errorCount   int // functions send error every time > 0
	loginOkCount int
}

// isMock returns true if client is mock.
func isMock(c im.ImmuClient) bool {
	_, ok := c.(*mockImmuClient)
	return ok
}

// getOkCount
func getOkCount(c im.ImmuClient) int {
	mock, ok := c.(*mockImmuClient)
	if ok {
		return mock.getOkCount
	}
	return 0
}

// setOkCount sets the mock's setOkCount, if mock exists.
func _(_ im.ImmuClient) int {
	mock, ok := immuLedger.client.(*mockImmuClient)
	if ok {
		return mock.setOkCount
	}
	return 0
}

// errorCount
func errorCount(c im.ImmuClient) int {
	mock, ok := c.(*mockImmuClient)
	if ok {
		return mock.errorCount
	}
	return 0
}

// setErrorCount
func setErrorCount(c im.ImmuClient, count int) {
	mock, ok := c.(*mockImmuClient)
	if ok {
		mock.errorCount = count
	}
}

// Override the real immuclient.Set() function. Can be used to return also errors if needed
func (m *mockImmuClient) Set(_ context.Context, key []byte, value []byte) (*schema.TxMetadata, error) {
	if m.errorCount > 0 {
		m.errorCount--
		return nil, errors.New("mock error")
	}
	glog.V(2).Infoln("mock set called with key:", string(key))
	// store values
	store[keyType(key)] = value
	// Set test data to return. This is how the real data looks like
	var txData schema.TxMetadata
	txData.Id = 108
	txData.PrevAlh = []byte("E+\x1e\x85\x85 X\x1d\x87\x8a\x03\xb1\xf2\xb1\xf5\x9eh\xa2\xf2_5{1Ӎ\x03Bٵڳ\xd9")
	txData.Ts = 1614767958
	txData.EH = []byte("BA\xaab\x9a{Y\xa4\xad\xd9\xee\xa4fn^^Q\x14d\x87k4%\xdcލC\xd6Ԁ\xc7(")
	txData.BlTxId = 107
	txData.BlRoot = []byte("q\xb7(<U]\xba\xad\x8b\xf1\x1cB\x83E\xe6`\xf9\xc3\x12\xe9y\x05\xf9+[\xfawS\xab\xa0\x92I")
	m.setOkCount++
	return &txData, nil
}

// Override the real immuclient.Get() function. Can be used to return also errors if needed
func (m *mockImmuClient) Get(_ context.Context, key []byte) (*schema.Entry, error) {
	if m.errorCount > 0 {
		m.errorCount--
		return nil, errors.New("mock error")
	}
	glog.V(2).Infoln("mock get called with key:", string(key))
	// Set test data to return. This is how the real data looks like
	var entryData schema.Entry
	entryData.Tx = 117
	entryData.Key = key
	entryData.Value = store[keyType(key)]
	m.getOkCount++
	return &entryData, nil
}

func (m *mockImmuClient) Login(
	_ context.Context,
	user []byte,
	pass []byte,
) (*schema.LoginResponse, error) {
	if m.errorCount > 0 {
		m.errorCount--
		return nil, errors.New("mock error")
	}
	glog.V(1).Infof("------------ immu mock login (user:%s/pwd:%s)",
		string(user), string(pass))
	lr := &schema.LoginResponse{Token: "MOCK_TOKEN"}
	m.loginOkCount++
	return lr, nil
}

func (m *mockImmuClient) Logout(_ context.Context) error {
	if m.errorCount > 0 {
		m.errorCount--
		return errors.New("mock error")
	}
	glog.V(1).Infoln("========= calling logout from mock db")
	return nil
}

// rmTokenDir removes empty token dir so simulate same behavior as the immuDB.
// The reason why we have to have it in mock is the bug in immuDB Logout
// function which is depending the existing of the directory.
func _() {
	hd, _ := os.UserHomeDir()
	fp := filepath.Join(hd, "token")
	err := os.Remove(fp)
	if err != nil {
		glog.Errorln("remove token:", err)
	}
	glog.V(12).Infoln("token path removed:", fp)
}
