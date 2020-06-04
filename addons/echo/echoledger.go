// Package echo is a ledger addon to offer implementation which writes a log
// about all of the ledger read/write.
package echo

import (
	"fmt"
	"sync"

	"github.com/optechlab/findy-go/pool"
)

const echoName = "FINDY_ECHO_LEDGER"

type ledger struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func (m *ledger) Close() {
	fmt.Println("Closing Echo ledger")
	resetMemory()
}

func (m *ledger) Open(name string) bool {
	fmt.Println("Opening Echo ledger")
	resetMemory()
	return name == echoName
}

func (m *ledger) Write(ID, data string) error {
	m.mem.Lock()
	defer m.mem.Unlock()
	fmt.Printf("Ledger WRITE [%s] <- (%s)", ID, data)
	m.mem.ory[ID] = data
	return nil
}

func (m *ledger) Read(ID string) (name string, value string, err error) {
	m.mem.RLock()
	defer m.mem.RUnlock()
	fmt.Printf("Ledger READ [%s] -> (%s)", ID, m.mem.ory[ID])
	return ID, m.mem.ory[ID], nil
}

var impl = &ledger{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

func init() {
	pool.RegisterPlugin(echoName, impl)
}

func resetMemory() {
	impl.mem.ory = make(map[string]string)
}
