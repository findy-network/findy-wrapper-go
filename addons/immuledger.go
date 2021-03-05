package addons

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"google.golang.org/grpc/metadata"
)

const immuLedgerName = "FINDY_IMMUDB_LEDGER"
const immuMockLedgerName = "FINDY_MOCK_IMMUDB_LEDGER"

type myImmuClient immuclient.ImmuClient

var client myImmuClient

type ImmuCfg struct {
	URL      string `json:"url"`
	Port     int    `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`

	immuclient.Options
}

// immu is a ledger addon which implements reading / writing data directly to the ImmuDB.
// It writes ledger data to memory before returning it and if it's stored in memory it serves
// the data from there instead of fetching it from the ImmuDB
type immu struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func tryGetOpions() (cfg *ImmuCfg) {
	// get credentials from env if available
	cfg = &ImmuCfg{
		URL:      os.Getenv("ImmuUrl"),
		Port:     err2.Int.Try(strconv.Atoi(os.Getenv("ImmuPort"))),
		UserName: os.Getenv("ImmuUsrName"),
		Password: os.Getenv("ImmuPasswd"),
	}
	assert.D.True(cfg.URL != "", "database URL is needed")
	assert.D.True(cfg.Port != 0, "port cannot be 0")
	assert.D.True(cfg.UserName != "", "user name cannot be empty")
	assert.D.True(cfg.Password != "", "password cannot be empty")
	cfg.Options = immuclient.Options{
		Address:         cfg.URL,
		Port:            cfg.Port,
		Auth:            true,
		CurrentDatabase: "defaultdb",
	}
	return cfg
}

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

// This is needed because of testing for clearing the memCache
func (i *immu) ResetMemCache() {
	resetImmuLedger()
}

func (i *immu) Close() {
	defer err2.Catch(func(err error) {
		glog.Errorf("error immu db ledger addon %v", err)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Logout(ctx)
	client = nil
	// resetImmuLedger()
}

func (i *immu) Open(name string) bool {
	// why this is reseted here? for test? should we load it from the DB at startup?
	resetImmuLedger()
	if name == immuMockLedgerName {
		// connection is done already, Mock is 'open'
		return true
	}
	Cfg = tryGetOpions()
	connectToImmu()

	return name == immuLedgerName
}

func (i *immu) Write(ID, data string) (err error) {
	defer err2.Return(&err)

	// is immuDB thread safe? Why this is undocumented method?

	// store the data to the memory cache in all cases
	i.mem.Lock()
	i.mem.ory[ID] = data
	i.mem.Unlock()

	// todo: extact this to own function later for retries, etc.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, err := client.Set(ctx, []byte(ID), []byte(data))
	err2.Check(err)
	fmt.Printf("Immuledger: Successfully committed key \"%s\" at tx %d\n", []byte(ID), tx.Id)
	// fmt.Println("Immuledger: tx ", tx)

	return nil
}

func (i *immu) Read(ID string) (name string, value string, err error) {
	// chekck if we have data in mem cache
	// for testing purposes, you might want to disable this temporarily when testing
	// to varify that reading from ImmuDB is succesful
	i.mem.RLock()
	if item, ok := i.mem.ory[ID]; ok {
		fmt.Printf("Immuledger: Successfully retrieved entry for key %s from memcache\n", ID)
		// data can be found from memcache, return it
		defer i.mem.RUnlock()
		return ID, item, nil
	}
	i.mem.RUnlock()

	// todo: extract to function to handle errors and retries
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dataFromImmu, err := client.Get(ctx, []byte(ID))
	if err != nil {
		fmt.Printf("Immuledger: Getting key \"%s\" from ImmuDB failed. Error \"%s\"", ID, err)
		err2.Check(err)
	}
	fmt.Printf("Immuledger: Successfully retrieved entry for key %s\n", dataFromImmu.Key)
	// fmt.Println("Immuledger: immudata ", dataFromImmu)

	i.mem.Lock()
	defer i.mem.Unlock()
	i.mem.ory[ID] = string(dataFromImmu.Value) // store the data to the memory cache
	return ID, string(dataFromImmu.Value), nil
}

var immuLedger = &immu{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

var Cfg = &ImmuCfg{
	URL:      "localhost",
	Port:     3322,
	UserName: "immudb",
	Password: "immudb",
}

func init() {
	pool.RegisterPlugin(immuLedgerName, immuLedger)
}

func resetImmuLedger() {
	immuLedger.mem.ory = make(map[string]string)
}
