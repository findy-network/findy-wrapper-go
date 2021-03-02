package addons

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/lainio/err2"
	"google.golang.org/grpc/metadata"
)

const immuLedgerName = "FINDY_IMMUDB_LEDGER"

var client immuclient.ImmuClient
var ctx context.Context

// immu is a ledger addon which implements reading / writing data directly to the ImmuDB.
// It writes ledger data to memory before returning it and if it's stored in memory it serves
// the data from there instead of fetching it from the ImmuDB
type immu struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

// clean possible line breaks and tabs from the data
func CleanDataString(data string) string {
	re := regexp.MustCompile(`\n`)
	data = re.ReplaceAllString(data, "")
	re = regexp.MustCompile(`\t`)
	data = re.ReplaceAllString(data, "")
	return data
}

func ConnectToImmu() (immuclient.ImmuClient, context.Context, error) {
	// get credentials from env
	immuURL := os.Getenv("ImmuUrl")
	immuPortString := os.Getenv("ImmuPort")
	userName := os.Getenv("ImmuUsrName")
	password := os.Getenv("ImmuPasswd")
	immuPort, err := strconv.Atoi(immuPortString)
	if err != nil {
		fmt.Println("Immuledger: Converting port to int failed", err)
		return nil, nil, err
	}
	// set connection options
	var options immuclient.Options
	options.Address = immuURL
	options.Port = immuPort
	options.Auth = true
	options.CurrentDatabase = "defaultdb"

	// connect to ImmuDB
	if client == nil {
		var err error
		client, err = immuclient.NewImmuClient(immuclient.DefaultOptions().WithAddress(immuURL).WithAuth(true))
		if err != nil {
			fmt.Println("Immuledger: Connect to ImmuDB failed", err)
			return nil, nil, err
		}
		ctx = context.Background()
		lr, err := client.Login(ctx, []byte(userName), []byte(password))
		if err != nil {
			fmt.Println("Immuledger: Login failed", err)
			return nil, nil, err
		}
		// immudb provides multidatabase capabilities.
		// token is used not only for authentication, but also to route calls to the correct database
		md := metadata.Pairs("authorization", lr.Token)
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		fmt.Println("Immuledger: Connection to ImmuDB is OK")
	}
	return client, ctx, nil
}

func (i *immu) Close() {
	client = nil
	ctx = nil
	// resetImmuLedger()
}

func (i *immu) Open(name string) bool {
	resetImmuLedger()
	return name == immuLedgerName
}

func (i *immu) Write(ID, data string) (err error) {
	data = CleanDataString(data) // do some data cleaning if needed

	// connect to ImmuDB
	client, ctx, err := ConnectToImmu()
	if err != nil {
		fmt.Println("Immuledger: Connect failed for writing", err)
		err2.Check(err)
	}
	tx, err := client.Set(ctx, []byte(ID), []byte(data))
	if err != nil {
		fmt.Printf("Immuledger: Storing data to immuDb failed for key\"%s\". Error \"%s\"", ID, err)
		err2.Check(err)
	}
	fmt.Printf("Immuledger: Successfully committed key \"%s\" at tx %d\n", []byte(ID), tx.Id)

	// store the data to the memory cache
	defer err2.Return(&err)
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
	client, ctx, err := ConnectToImmu()
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
