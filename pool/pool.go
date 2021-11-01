/*
Package pool is corresponding Go package for libindy's pool namespace. We
suggest that you read indy SDK documentation for more information.

Pool offers functions to manage ledger pools. The Go wrapper has extended pool
system which means that it supports ledger plugins to implement own
persistent storage. These additional ledger can be used all at the same time:

	select {
	case r = <-pool.OpenLedger(poolName, "FINDY_MEM_LEDGER"):
		if r.Err() != nil {
			log.Panicln("Cannot open pool")
		}
	case <-time.After(maxTimeout):
		log.Panicln("Timeout exceeded")
	}
	h := r.Handle()

Current implementation offers a memory ledger for tests and simple cache. The
interface to use them is an extension to OpenLedger function which takes
multiple ledger pool names at once.
*/
package pool

import (
	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/findy-network/findy-wrapper-go/internal/c2go"
	"github.com/findy-network/findy-wrapper-go/internal/ctx"
	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
)

type readChan chan readInfo

type readInfo struct {
	id     string
	result string
	err    error
}

var (
	// Counter for plugin handles. It gives unique handles for plugins.
	// The handles are negative numbers: -1, -2, ..
	pluginHandles = -1

	// Currently opened plugins
	openPlugins = make(map[int]plugin.Ledger)
)

// RegisterPlugin is interface for ledger plugins to register them selves.
func RegisterPlugin(name string, plugin plugin.Ledger) {
	registeredPlugins[name] = plugin
}

// CreateConfig is indy SDK wrapper to create ledger pool configuration. See
// more information from indy_pool_create_config().
func CreateConfig(name string, config Config) ctx.Channel {
	return c2go.PoolCreateConfig(name, dto.ToJSON(config))
}

// List is indy SDK wrapper to list current pool configurations on the host. See
// more information from indy_pool_list().
func List() ctx.Channel {
	return c2go.PoolList()
}

// ListPlugins is ledger plugin system function. With the function caller can
// have list of all currently installed and activated ledger plugins in the
// compilation.
func ListPlugins() []string {
	names := make([]string, 0, len(registeredPlugins))
	for name := range registeredPlugins {
		names = append(names, name)
	}
	return names
}

// OpenLedger opens all ledger types given. The original indy SDK function takes
// only one pool configuration name as an argument. This wrapper takes many pool
// plugin name-argument pairs. ListPlugins() returns all of the available ones.
func OpenLedger(names ...string) ctx.Channel {
	_, startsWithPluginName := registeredPlugins[names[0]]
	legacy := len(names) == 1

	// only one argument is given, as legacy mode was
	if legacy {
		// default is that incoming argument is Indy pool name only
		pluginName := "FINDY_LEDGER"
		pluginArg := names[0]

		if startsWithPluginName {
			pluginName = names[0]
			pluginArg = ""
		}
		names = make([]string, 2)
		names[0] = pluginName
		names[1] = pluginArg
	}

	for i := 0; i < len(names); i += 2 {
		name := names[i]
		extra := names[i+1]
		if r, ok := registeredPlugins[name]; ok && r.Open(name, extra) {
			openPlugins[pluginHandles] = r
			pluginHandles--
		}
	}
	return makeHandleResult(pluginHandles + 1)
}

// CloseLedger is exchanged indy SDK wrapper function. It closes all currently
// opened ledger plugins and the actual indy ledger if it is open as well.
func CloseLedger(handle int) ctx.Channel {
	for _, ledger := range openPlugins {
		ledger.Close()
	}
	pluginHandles = -1
	openPlugins = make(map[int]plugin.Ledger)
	if handle > 0 {
		return c2go.PoolCloseLedger(handle)
	}
	return makeHandleResult(0)
}

// SetProtocolVersion is indy SDK wrapper. It sets the used protocol version. In
// most cases it is 2. See more information from indy_set_protocol_version().
func SetProtocolVersion(version uint64) ctx.Channel {
	return c2go.PoolSetProtocolVersion(version)
}

// IsIndyLedgerOpen return true if indy ledger is open as well. It is used to
// coordinate ledger transactions.
func IsIndyLedgerOpen(handle int) bool {
	return handle > 0
}

// Write writes data to all of the plugin ledgers. Note! The original indy
// ledger is not one of the plugin ledgers, at least not yet.
func Write(tx plugin.TxInfo, ID, data string) {
	for _, ledger := range openPlugins {
		l := ledger
		go func() {
			if err := l.Write(tx, ID, data); err != nil {
				glog.Errorln("-- writing err ledger:", tx, ID, "data:\n", data)
				glog.Errorln("-- error:", err)
			}
		}()
	}
}

// Read reads data from all of the plugin ledgers. Note! The original indy
// ledger is not one of the plugin ledgers, at least not yet.
func Read(tx plugin.TxInfo, ID string) (string, string, error) {
	switch len(openPlugins) {
	case 0:
		assert.D.True(false, "no plugins open")
	case 1:
		return openPlugins[-1].Read(tx, ID)
	case 2:
		return readFrom2(tx, ID)
	default:
		assert.D.True(false, "amount of open plugins is not supported")
	}
	return "", "", nil
}

func readFrom2(tx plugin.TxInfo, ID string) (id string, val string, err error) {
	defer err2.Annotate("reading cached ledger", &err)

	const (
		indyLedger  = -1
		cacheLedger = -2
	)
	var (
		result    string
		readCount int
	)

	ch1 := asyncRead(indyLedger, tx, ID)
	ch2 := asyncRead(cacheLedger, tx, ID)

loop:
	for {
		select {
		case r1 := <-ch1:
			if r1.err != nil && r1.err != plugin.ErrNotExist {
				err2.Check(r1.err)
			}
			readCount++
			glog.V(5).Infof("---- %d. winner -1 ----", readCount)
			result = r1.result

			// Currently first plugin is the Indy ledger, if we are
			// here, we must write data to cache ledger
			if readCount >= 2 && r1.result != "" {
				glog.V(5).Infoln("--- update cache plugin:", r1.id, r1.result)
				tmpTx := tx
				tx.Update = true
				err := openPlugins[cacheLedger].Write(tmpTx, ID, r1.result)
				if err != nil {
					glog.Errorln("error cache update", err)
				}
			}
			break loop

		case r2 := <-ch2:
			if r2.err != nil && r2.err != plugin.ErrNotExist {
				err2.Check(r2.err)
			}
			readCount++
			glog.V(5).Infof("---- %d. winner -2 ----", readCount)
			result = r2.result
			if r2.result == "" {
				glog.V(5).Infoln("--- NO CACHE HIT:", ID, readCount)
				continue loop
			}
			break loop
		}
	}
	return ID, result, nil
}

func asyncRead(i int, tx plugin.TxInfo, ID string) readChan {
	ch := make(readChan)
	go func() {
		name, value, err := openPlugins[i].Read(tx, ID)
		if err != nil {
			glog.Errorf("error in value: %s, ledger reading: %s", name, err)
		}
		ch <- readInfo{
			id:     name,
			result: value,
		}
	}()
	return ch
}

// registeredPlugins keeps track of the all of leger plugins installed in the
// compilation.
var registeredPlugins = make(map[string]plugin.Ledger)

// makeHandleResult makes and returns context channel to return for the caller.
// It is helper for cases where indy wrapping system is not used i.e. indy
// callback is not called.
func makeHandleResult(h int) ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.Push()
	go func() {
		r := dto.Result{
			Data: dto.Data{
				Handle: h,
			},
		}
		c := ctx.CmdContext.Pop(cmdHandle, r)
		c <- r
	}()
	return ch
}
