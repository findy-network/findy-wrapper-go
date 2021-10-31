package addons

import (
	"sync"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/findy-network/findy-wrapper-go/pool"
)

const memName = "FINDY_MEM_LEDGER"

// Mem is a ledger addon which implements transient ledger. It writes
// ledger data to memory and reads it from there. It's convenient for unit test
// and some development cases.
type Mem struct {
	Mem struct {
		sync.RWMutex
		Ory map[string]string
	}
}

func (m *Mem) Close() {
	resetMem()
}

func (m *Mem) Open(name ...string) bool {
	resetMem()
	return name[0] == memName
}

func (m *Mem) Write(_ plugin.TxInfo, ID, data string) error {
	m.Mem.Lock()
	defer m.Mem.Unlock()
	m.Mem.Ory[ID] = data
	return nil
}

func (m *Mem) Read(_ plugin.TxInfo, ID string) (name string, value string, err error) {
	m.Mem.RLock()
	defer m.Mem.RUnlock()
	return ID, m.Mem.Ory[ID], nil
}

var memLedger = &Mem{Mem: struct {
	sync.RWMutex
	Ory map[string]string
}{}}

func init() {
	pool.RegisterPlugin(memName, memLedger)
}

func resetMem() {
	memLedger.Mem.Ory = make(map[string]string)
}
