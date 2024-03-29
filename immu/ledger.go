package immu

import (
	"context"
	"math"
	"sync"
	"time"

	im "github.com/codenotary/immudb/pkg/client"
	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"google.golang.org/grpc/metadata"
)

type myImmuClient im.ImmuClient

type immu struct {
	cache  Mem
	client myImmuClient
	token  string
	cfg    *Cfg
}

func (i *immu) Close() {
	defer err2.Catch(err2.Err(func(err error) {
		glog.Errorf("error immu db ledger addon Close(): %v", err)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = i.buildCtx(ctx)
	try.To(i.client.Logout(ctx))

	i.client = nil
}

func (i *immu) Open(name ...string) bool {
	defer err2.Catch(err2.Err(func(err error) {
		glog.Errorf("error immu db ledger addon Open(): %v", err)
	}))
	i.ResetMemCache() // for tests at the moment

	cfg := NewImmuCfg(name[0])
	c, token := try.To2(cfg.Connect())

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

func (i *immu) Write(tx plugin.TxInfo, ID, data string) (err error) {
	defer err2.Handle(&err, func(err error) error {
		glog.Errorln("write error:", err)
		err = nil // suspend error for now and retry
		go i.writeRetry(ID, data)
		return err
	})

	_ = i.cache.Write(tx, ID, data)
	try.To(i.oneWrite(ID, data))
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
	defer err2.Handle(&err, func(err error) error {
		glog.Errorf("retry db write: %v", err)
		return err
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = i.buildCtx(ctx)
	try.To1(i.client.Set(ctx, []byte(ID), []byte(data)))
	return nil
}

func (i *immu) Read(tx plugin.TxInfo, ID string) (name string, value string, err error) {
	defer err2.Handle(&err)

	if _, value, err = i.cache.Read(tx, ID); err == nil && value != "" {
		glog.V(1).Info("----- cache hit")
		return ID, value, err
	}
	if err != nil {
		glog.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = i.buildCtx(ctx)
	dataFromImmu := try.To1(i.client.Get(ctx, []byte(ID)))

	_ = i.cache.Write(tx, ID, string(dataFromImmu.Value))
	return ID, string(dataFromImmu.Value), nil
}

func (i *immu) login() (err error) {
	defer err2.Handle(&err)
	i.token = try.To1(i.cfg.login(i.client))
	return nil
}

func (i *immu) ResetMemCache() {
	glog.V(1).Infof("------------ reset cache (%d)", len(i.cache.Mem.Ory))
	i.cache.Mem.Lock()
	i.cache.Mem.Ory = make(map[string]string)
	i.cache.Mem.Unlock()
}

var _ = Mem{Mem: struct {
	sync.RWMutex
	Ory map[string]string
}{}}

var immuLedgerImpl = &immu{cache: Mem{Mem: struct {
	sync.RWMutex
	Ory map[string]string
}{}}}

func init() {
	//pool.RegisterPlugin(immuLedgerNameImpl, immuLedger)
}
