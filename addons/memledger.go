package addons

import (
	"sync"

	"github.com/findy-network/findy-wrapper-go/pool"
)

const memName = "FINDY_MEM_LEDGER"

// mem is a ledger addon which implements transient ledger. It writes
// ledger data to memory and reads it from there. It's convenient for unit test
// and some development cases.
type mem struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func (m *mem) Close() {
	resetMem()
}

func (m *mem) Open(name string) bool {
	resetMem()
	return name == memName
}

func (m *mem) Write(ID, data string) error {
	m.mem.Lock()
	defer m.mem.Unlock()
	m.mem.ory[ID] = data
	return nil
}

func (m *mem) Read(ID string) (name string, value string, err error) {
	m.mem.RLock()
	defer m.mem.RUnlock()
	return ID, m.mem.ory[ID], nil
}

var memLedger = &mem{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

func init() {
	pool.RegisterPlugin(memName, memLedger)
}

func resetMem() {
	memLedger.mem.ory = make(map[string]string)
}
