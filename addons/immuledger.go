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
	"github.com/lainio/err2"
	"google.golang.org/grpc/metadata"
)

const immuLedgerName = "FINDY_IMMUDB_LEDGER"

type myImmuClient immuclient.ImmuClient

var client myImmuClient

// immu is a ledger addon which implements reading / writing data directly to the ImmuDB.
// It writes ledger data to memory before returning it and if it's stored in memory it serves
// the data from there instead of fetching it from the ImmuDB
type immu struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func connectToImmu() (_ immuclient.ImmuClient, _ context.Context, err error) {
	defer err2.Return(&err)

	// get credentials from env
	immuURL := os.Getenv("ImmuUrl")
	immuPortString := os.Getenv("ImmuPort")
	userName := os.Getenv("ImmuUsrName")
	password := os.Getenv("ImmuPasswd")
	immuPort := err2.Int.Try(strconv.Atoi(immuPortString))
	// set connection options
	options := &immuclient.Options{
		Address:         immuURL,
		Port:            immuPort,
		Auth:            true,
		CurrentDatabase: "defaultdb",
	}

	// connect to ImmuDB
	if client == nil {
		client, err = immuclient.NewImmuClient(options)
		err2.Check(err)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		lr, err := client.Login(ctx, []byte(userName), []byte(password))
		err2.Check(err)
		// immudb provides multidatabase capabilities.
		// token is used not only for authentication, but also to route calls to the correct database
		md := metadata.Pairs("authorization", lr.Token)
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		fmt.Println("Immuledger: Connection to ImmuDB is OK")
	}
	return client, nil, nil
}

// This is needed because of testing for clearing the memCache
func (i *immu) ResetMemCache() {
	resetImmuLedger()
}

func (i *immu) Close() {
	client = nil
	// resetImmuLedger()
}

func (i *immu) Open(name string) bool {
	// why this is reseted here? should we load it from the DB at startup?
	resetImmuLedger()

	return name == immuLedgerName
}

func (i *immu) Write(ID, data string) (err error) {
	defer err2.Return(&err)

	// connect to ImmuDB
	client, ctx, err := connectToImmu()
	err2.Check(err)
	// is immuDB thread safe? Why this is undocumented method?
	// what happens to db connection during the error
	tx, err := client.Set(ctx, []byte(ID), []byte(data))
	err2.Check(err)
	fmt.Printf("Immuledger: Successfully committed key \"%s\" at tx %d\n", []byte(ID), tx.Id)
	// fmt.Println("Immuledger: tx ", tx)

	// store the data to the memory cache
	i.mem.Lock()
	defer i.mem.Unlock()
	i.mem.ory[ID] = data

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

	// connect to ImmuDB
	client, ctx, err := connectToImmu()
	if err != nil {
		fmt.Println("Immuledger: Connect failed for reading", err)
		err2.Check(err)
	}
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

func init() {
	// env settings. Probably these will be set where this addon is called
	// So these lines will be removed from here eventually
	os.Setenv("ImmuUrl", "localhost")
	os.Setenv("ImmuPort", "3322")
	os.Setenv("ImmuUsrName", "immudb")
	os.Setenv("ImmuPasswd", "immudb")
	pool.RegisterPlugin(immuLedgerName, immuLedger)
}

func resetImmuLedger() {
	immuLedger.mem.ory = make(map[string]string)
}
