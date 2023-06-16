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
	"fmt"
	"strings"
	"sync"

	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/findy-network/findy-wrapper-go/internal/c2go"
	"github.com/findy-network/findy-wrapper-go/internal/ctx"
	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
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
	// first round of checks, if caller cannot yet use variadic function
	if len(names) == 1 {
		names = ConvertPluginArgs(names[0])
	}

	names = BuildLegacyPluginArgs(names)

	for i := 0; i < len(names); i += 2 {
		name := names[i]
		extra := names[i+1]
		lstr := fmt.Sprintf("open plugin:%s(%s)", name, extra)
		if r, ok := registeredPlugins[name]; ok && r.Open(extra) {
			openPlugins[pluginHandles] = r
			pluginHandles--
			lstr += " ==> OK"
		} else {
			lstr += " ==> ERR"
		}
		glog.V(1).Infoln(lstr)
	}
	return makeHandleResult(pluginHandles + 1)
}

func BuildLegacyPluginArgs(names []string) (ns []string) {
	_, startsWithPluginName := registeredPlugins[names[0]]
	legacy := len(names) == 1 && !startsWithPluginName

	// only one argument is given, as legacy mode was
	if !legacy {
		glog.V(3).Infoln("pluginName is NOT legacy")
		return names
	}
	glog.V(3).Infoln("pluginName is legacy")
	// default is that incoming argument is Indy pool name only
	pluginName := "FINDY_LEDGER"
	pluginArg := names[0]

	if startsWithPluginName {
		pluginName = names[0]
		pluginArg = ""
	}
	ns = make([]string, 2)
	ns[0] = pluginName
	ns[1] = pluginArg
	return ns
}

// ConvertPluginArgs converts pool name to list of config pairs if it's possible
func ConvertPluginArgs(poolName string) []string {
	glog.V(3).Infoln("convert plugin args for:", poolName)

	pools := strings.Split(poolName, ",")
	pluginsLen := len(pools)
	if pluginsLen == 1 {
		_, startsWithPluginName := registeredPlugins[poolName]
		if !startsWithPluginName {
			return []string{poolName}
		}

		pluginsLen++
		pools = append(pools, "")
	}
	glog.V(1).Infof("Using env var defined %d ledger plugin(s)", pluginsLen)

	poolNames := make([]string, pluginsLen)
	for i := 0; i < pluginsLen; i += 2 {
		poolNames[i] = pools[i]
		poolNames[i+1] = pools[i+1]
	}
	if pluginsLen >= 2 && glog.V(1) {
		glog.Infoln("Using two ledger plugins")
	}

	return poolNames
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

// Write writes data to all of the plugin ledgers, and waits their results.
func Write(tx plugin.TxInfo, ID, data string) (err error) {
	var wg sync.WaitGroup
	for _, ledger := range openPlugins {
		l := ledger
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer err2.Catch(func(er error) {
				glog.Errorln("-- writing err ledger:", tx, ID, "data:\n", data)
				glog.Errorln("-- error:", er)
				if err != nil {
					// if there was previous error, get them both
					err = fmt.Errorf("plugin write error: %w: %w", err, er)
				} else {
					// other report just current error (er)
					err = fmt.Errorf("plugin write error: %w", er)
				}
			})
			try.To(l.Write(tx, ID, data))
		}()
	}
	wg.Wait()
	return err
}

// Read reads data from all of the plugin ledgers.
func Read(tx plugin.TxInfo, ID string) (string, string, error) {
	switch len(openPlugins) {
	case 0:
		assert.That(false, "no plugins open")
	case 1:
		return openPlugins[-1].Read(tx, ID)
	case 2:
		return readFrom2(tx, ID)
	default:
		assert.That(false, "amount of open plugins is not supported")
	}
	return "", "", nil
}

func readFrom2(tx plugin.TxInfo, ID string) (id string, val string, err error) {
	defer err2.Handle(&err, "reading cached ledger")

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
			exist := !try.Is(r1.err, plugin.ErrNotExist)

			readCount++
			glog.V(5).Infof("---- %d. winner -1 (exist=%v) ----",
				readCount, exist)
			result = r1.result

			// Currently first plugin is the Indy ledger, if we are
			// here, we must write data to cache ledger
			if readCount >= 2 && exist {
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
			notExist := try.Is(r2.err, plugin.ErrNotExist)

			readCount++
			glog.V(5).Infof("---- %d. winner -2 (notExist=%v, result=%s) ----",
				readCount, notExist, r2.result)
			result = r2.result

			if notExist {
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
			err:    err,
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
