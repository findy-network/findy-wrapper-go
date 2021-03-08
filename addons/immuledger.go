package addons

import (
	"context"
	"sync"
	"time"

	im "github.com/codenotary/immudb/pkg/client"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"google.golang.org/grpc/metadata"
)

const immuLedgerName = "FINDY_IMMUDB_LEDGER"

type myImmuClient im.ImmuClient

type immu struct {
	cache  mem
	client myImmuClient
	token  string
}

func (i *immu) Close() {
	defer err2.Catch(func(err error) {
		glog.Errorf("error immu db ledger addon Close(): %v", err)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	i.client.Logout(ctx)
	i.client = nil
}

func (i *immu) Open(name string) bool {
	defer err2.Catch(func(err error) {
		glog.Errorf("error immu db ledger addon Open(): %v", err)
	})
	// why this is reseted here? for test? should we load it from the DB at startup?
	i.ResetMemCache()

	cfg := NewImmuCfg(name)
	c, token, err := cfg.Connect()
	err2.Check(err)

	i.client = c
	i.token = token
	return true
}

func (i *immu) buildCtx(context context.Context) context.Context {
	// immudb provides multidatabase capabilities.
	// token is used not only for authentication, but also to route calls to the correct database
	md := metadata.Pairs("authorization", i.token)
	ctx := metadata.NewOutgoingContext(context, md)
	return ctx
}

func (i *immu) Write(ID, data string) (err error) {
	defer err2.Return(&err)

	_ = i.cache.Write(ID, data)

	// todo: extact this to own function later for retries, etc.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = i.buildCtx(ctx)
	_, err = i.client.Set(ctx, []byte(ID), []byte(data))
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
	ctx = i.buildCtx(ctx)
	dataFromImmu, err := i.client.Get(ctx, []byte(ID))
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

func init() {
	pool.RegisterPlugin(immuLedgerName, immuLedger)
}
