package addons

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
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

func tryOptions() (cfg *ImmuCfg) {
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

type immu struct {
	cache mem
}

func (i *immu) Close() {
	defer err2.Catch(func(err error) {
		glog.Errorf("error immu db ledger addon %v", err)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Logout(ctx)
	client = nil
}

func (i *immu) Open(name string) bool {
	// why this is reseted here? for test? should we load it from the DB at startup?
	i.ResetMemCache()

	if name == immuMockLedgerName {
		// connection is done already, Mock is 'open'
		return true
	}
	Cfg = tryOptions()
	connectToImmu()

	return name == immuLedgerName
}

func (i *immu) Write(ID, data string) (err error) {
	defer err2.Return(&err)

	_ = i.cache.Write(ID, data)

	// todo: extact this to own function later for retries, etc.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.Set(ctx, []byte(ID), []byte(data))
	err2.Check(err)
	//fmt.Printf("Immuledger: Successfully committed key \"%s\" at tx %d\n", []byte(ID), tx.Id)
	// fmt.Println("Immuledger: tx ", tx)

	return nil
}

func (i *immu) Read(ID string) (name string, value string, err error) {
	defer err2.Return(&err)

	// chekck if we have data in mem cache
	if _, value, err = i.cache.Read(ID); err == nil && value != "" {
		glog.V(1).Info("----- cache hit")
		return ID, value, err
	}
	if err != nil {
		glog.Error(err)
	}

	// todo: extract to function to handle errors and retries
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dataFromImmu, err := client.Get(ctx, []byte(ID))
	err2.Check(err)
	//fmt.Printf("Immuledger: Successfully retrieved entry for key %s\n", dataFromImmu.Key)

	_ = i.cache.Write(ID, string(dataFromImmu.Value))
	return ID, string(dataFromImmu.Value), nil
}

func (i *immu) ResetMemCache() {
	glog.V(1).Infof("------------ reset cache (%d)", len(i.cache.mem.ory))
	i.cache.mem.Lock()
	i.cache.mem.ory = make(map[string]string)
	i.cache.mem.Unlock()
}

var immuMemLedger = mem{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

var immuLedger = &immu{cache: mem{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}}

var Cfg = &ImmuCfg{
	URL:      "localhost",
	Port:     3322,
	UserName: "immudb",
	Password: "immudb",
}

func init() {
	pool.RegisterPlugin(immuLedgerName, immuLedger)
}
