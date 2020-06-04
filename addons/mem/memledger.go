// Package mem is a ledger addon which implements transient ledger. It writes
// ledger data to memory and reads it from there. It's convenient for unit test
// and some development cases.
package mem

import (
	"sync"

	"github.com/optechlab/findy-go/pool"
)

const addonName = "FINDY_MEM_LEDGER"

type addon struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func (m *addon) Close() {
	resetMemory()
}

func (m *addon) Open(name string) bool {
	resetMemory()
	return name == addonName
}

func (m *addon) Write(ID, data string) error {
	m.mem.Lock()
	defer m.mem.Unlock()
	m.mem.ory[ID] = data
	return nil
}

func (m *addon) Read(ID string) (name string, value string, err error) {
	m.mem.RLock()
	defer m.mem.RUnlock()
	return ID, m.mem.ory[ID], nil
}

var ledger = &addon{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

func init() {
	pool.RegisterPlugin(addonName, ledger)
}

func resetMemory() {
	ledger.mem.ory = make(map[string]string)
}
