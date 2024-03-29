// Package addons includes ledger addons.
package addons

import (
	"fmt"
	"sync"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/findy-network/findy-wrapper-go/pool"
)

const echoName = "FINDY_ECHO_LEDGER"

// echo offers implementation which writes a log
// about all of the ledger read/write.
type echo struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func (m *echo) Close() {
	fmt.Println("Closing Echo ledger")
	resetEcho()
}

func (m *echo) Open(_ ...string) bool {
	fmt.Println("Opening Echo ledger")
	resetEcho()
	return true
}

func (m *echo) Write(_ plugin.TxInfo, ID, data string) error {
	m.mem.Lock()
	defer m.mem.Unlock()
	fmt.Printf("Ledger WRITE [%s] <- (%s)", ID, data)
	m.mem.ory[ID] = data
	return nil
}

func (m *echo) Read(_ plugin.TxInfo, ID string) (name string, value string, err error) {
	m.mem.RLock()
	defer m.mem.RUnlock()
	fmt.Printf("Ledger READ [%s] -> (%s)", ID, m.mem.ory[ID])
	return ID, m.mem.ory[ID], nil
}

var impl = &echo{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

func init() {
	pool.RegisterPlugin(echoName, impl)
}

func resetEcho() {
	impl.mem.ory = make(map[string]string)
}
