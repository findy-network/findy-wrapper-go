package addons

import (
	"context"
	"math"
	"sync"
	"time"

	im "github.com/codenotary/immudb/pkg/client"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"google.golang.org/grpc/metadata"
)

type myImmuClient im.ImmuClient

type immu struct {
	cache  mem
	client myImmuClient
	token  string
	cfg    *ImmuCfg
}

func (i *immu) Close() {
	defer err2.Catch(func(err error) {
		glog.Errorf("error immu db ledger addon Close(): %v", err)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = i.buildCtx(ctx)
	err2.Check(i.client.Logout(ctx))

	i.client = nil
}

func (i *immu) Open(name string) bool {
	defer err2.Catch(func(err error) {
		glog.Errorf("error immu db ledger addon Open(): %v", err)
	})
	i.ResetMemCache() // for tests at the moment

	cfg := NewImmuCfg(name)
	c, token, err := cfg.Connect()
	err2.Check(err)

	i.cfg = cfg
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
	defer err2.Handle(&err, func() {
		glog.Errorln("write error:", err)
		err = nil // suspend error for now and retry
		go i.writeRetry(ID, data)
	})

	_ = i.cache.Write(ID, data)
	err2.Check(i.oneWrite(ID, data))
	return nil
}

func (i *immu) writeRetry(ID, data string) {
	success := false
	x := 2.0
	for round := 0.0; round < 12.0; round++ {
		v := time.Duration(math.Pow(x, round))
		time.Sleep(v * time.Second)
		if err := i.oneWrite(ID, data); err != nil {
			success = true
			glog.Info("succesful db write retry")
			break
		}
	}
	if !success {
		glog.Error("cannot write to DB giving up")
	}
}

func (i *immu) oneWrite(ID, data string) (err error) {
	defer err2.Handle(&err, func() {
		glog.Errorf("retry db write: %v", err)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = i.buildCtx(ctx)
	_, err = i.client.Set(ctx, []byte(ID), []byte(data))
	err2.Check(err)
	return nil
}

func (i *immu) Read(ID string) (name string, value string, err error) {
	defer err2.Return(&err)

	if _, value, err = i.cache.Read(ID); err == nil && value != "" {
		glog.V(1).Info("----- cache hit")
		return ID, value, err
	}
	if err != nil {
		glog.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = i.buildCtx(ctx)
	dataFromImmu, err := i.client.Get(ctx, []byte(ID))
	err2.Check(err)

	_ = i.cache.Write(ID, string(dataFromImmu.Value))
	return ID, string(dataFromImmu.Value), nil
}

func (i *immu) login() (err error) {
	defer err2.Return(&err)
	i.token = err2.String.Try(i.cfg.login(i.client))
	return nil
}

func (i *immu) ResetMemCache() {
	glog.V(1).Infof("------------ reset cache (%d)", len(i.cache.mem.ory))
	i.cache.mem.Lock()
	i.cache.mem.ory = make(map[string]string)
	i.cache.mem.Unlock()
}

var _ = mem{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

var immuLedgerImpl = &immu{cache: mem{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}}

func init() {
	//pool.RegisterPlugin(immuLedgerNameImpl, immuLedger)
}
