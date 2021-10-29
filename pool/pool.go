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
	"github.com/lainio/err2/assert"
)

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
// names where one can be an original indy pool and the rest can be any of the
// plugins compiled into the binary running. ListPlugins() returns all of the
// available ones.
func OpenLedger(names ...string) ctx.Channel {
	realName := ""
	for _, name := range names {
		if r, ok := registeredPlugins[name]; ok && r.Open(name) {
			openPlugins[pluginHandles] = r
			pluginHandles--
		} else {
			if realName != "" {
				glog.Warning(
					"trying open multiple real ledgers",
					realName, " and", name, "using", name)
			}
			realName = name
		}
	}
	if realName != "" {
		return c2go.PoolOpenLedger(realName)
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
		ledger := ledger
		go func() {
			if err := ledger.Write(tx, ID, data); err != nil {
				glog.Error("error in writing ledger:", err)
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
		var result string
		select {
		case r1 := <-asyncRead(-1, tx, ID):
			result = r1
		case r2 := <-asyncRead(-2, tx, ID):
			result = r2
		}
		return ID, result, nil
	default:
		assert.D.True(false, "not suppoted plugins open")
	}
	return "", "", nil
}

func asyncRead(i int, tx plugin.TxInfo, ID string) chan string {
	ch := make(chan string)
	go func() {
		name, value, err := openPlugins[i].Read(tx, ID)
		if err != nil {
			glog.Errorf("error in value: %s, ledger reading: %s", name, err)
		}
		ch <- value
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
