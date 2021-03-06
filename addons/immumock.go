package addons

import (
	"context"
	"fmt"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"google.golang.org/grpc/metadata"
)

var storedKey []byte
var storedValue []byte

// create immuClient mock
type mockImmuClient struct {
	immuclient.ImmuClient
}

// This is needed because of testing. Using this the immuclient library
// functions Set() and Get() can be overriden and there is no need
// to have connectivity towards real ImmuDB
func MockImmuClientForTesting(newImmuclient myImmuClient) {
	client = newImmuclient
}

// Override the real immuclient.Set() function. Can be used to return also errors if needed
func (m *mockImmuClient) Set(ctx context.Context, key []byte, value []byte) (*schema.TxMetadata, error) {
	glog.V(2).Infoln("mock set called with key:", string(key))
	// store values
	storedKey = key
	storedValue = value
	// Set test data to return. This is how the real data looks like
	var txData schema.TxMetadata
	txData.Id = 108
	txData.PrevAlh = []byte("E+\x1e\x85\x85 X\x1d\x87\x8a\x03\xb1\xf2\xb1\xf5\x9eh\xa2\xf2_5{1Ӎ\x03Bٵڳ\xd9")
	txData.Ts = 1614767958
	txData.EH = []byte("BA\xaab\x9a{Y\xa4\xad\xd9\xee\xa4fn^^Q\x14d\x87k4%\xdcލC\xd6Ԁ\xc7(")
	txData.BlTxId = 107
	txData.BlRoot = []byte("q\xb7(<U]\xba\xad\x8b\xf1\x1cB\x83E\xe6`\xf9\xc3\x12\xe9y\x05\xf9+[\xfawS\xab\xa0\x92I")
	return &txData, nil
}

// Override the real immuclient.Get() function. Can be used to return also errors if needed
func (m *mockImmuClient) Get(ctx context.Context, key []byte) (*schema.Entry, error) {
	glog.V(2).Infoln("mock get called with key:", string(key))
	// Set test data to return. This is how the real data looks like
	var entryData schema.Entry
	entryData.Tx = 117
	entryData.Key = storedKey
	entryData.Value = storedValue
	return &entryData, nil
}

//func NewImmuClient(options *Options) (c ImmuClient, err error) {
//}

func connectToImmu() (err error) {
	defer err2.Return(&err)

	assert.P.True(client != nil, "client connection cannot be already open")

	client, err = immuclient.NewImmuClient(&Cfg.Options)
	err2.Check(err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lr, err := client.Login(ctx, []byte(Cfg.UserName), []byte(Cfg.Password))
	err2.Check(err)
	// immudb provides multidatabase capabilities.
	// token is used not only for authentication, but also to route calls to the correct database
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	fmt.Println("Immuledger: Connection to ImmuDB is OK")
	return nil
}
