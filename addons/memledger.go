package addons

import (
	"strconv"
	"strings"
	"sync"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
)

const memName = "FINDY_MEM_LEDGER"

// Mem is a ledger addon which implements transient ledger. It writes
// ledger data to memory and reads it from there. It's convenient for unit test
// and some development cases.
type Mem struct {
	Mem struct {
		sync.Mutex
		Ory map[string]string
	}

	// Seq is seqNo in real Indy ledger, by this we get correct behaviour
	Seq struct {
		sync.Mutex
		No uint
	}

	cacheMode bool
}

func (m *Mem) Close() {
	m.resetMem()
}

func (m *Mem) Open(name ...string) bool {
	m.resetMem()
	m.IncSeqNo()

	if name[0] != "" {
		glog.V(3).Infoln("-- setting Cache Mode for mem plugin --")
		m.cacheMode = true
	}

	return true
}

func (m *Mem) Write(tx plugin.TxInfo, ID, data string) error {
	//	if !tx.Update && tx.TxType == plugin.TxTypeSchema {
	//		glog.V(1).Infoln("----------- debugging -------")
	//		return nil
	//	}

	m.Mem.Lock()
	defer m.Mem.Unlock()

	m.IncSeqNo()
	m.Mem.Ory[ID] = data
	return nil
}

func (m *Mem) Read(tx plugin.TxInfo, ID string) (name string, value string, err error) {
	m.Mem.Lock()
	defer m.Mem.Unlock()

	seqNo := m.SeqNo()
	curval := m.Mem.Ory[ID]

	// if reading first time ("null" exists in "seqNo:"), replace it with the
	// current seqNo. This mimics indy ledger behaviour
	if tx.TxType == plugin.TxTypeSchema && strings.Contains(curval, "null") {
		if m.cacheMode {
			return ID, "", plugin.ErrNotExist
		}
		repval := strings.Replace(curval, "null", strconv.Itoa(int(seqNo)), 1)
		m.Mem.Ory[ID] = repval
	}

	return ID, m.Mem.Ory[ID], nil
}

func (m *Mem) IncSeqNo() {
	m.Seq.Lock()
	defer m.Seq.Unlock()

	m.Seq.No++
}

func (m *Mem) SeqNo() uint {
	m.Seq.Lock()
	defer m.Seq.Unlock()

	return m.Seq.No
}

func (m *Mem) resetMem() {
	m.Mem.Lock()
	defer m.Mem.Unlock()

	glog.V(3).Infoln("memLedger reset mem")
	m.Mem.Ory = make(map[string]string)
}

var memLedger = &Mem{
	Mem: struct {
		sync.Mutex
		Ory map[string]string
	}{},
	Seq: struct {
		sync.Mutex
		No uint
	}{
		No: 4, // Just installed empty Indy ledger starts about from here
	},
}

func init() {
	pool.RegisterPlugin(memName, memLedger)
}
