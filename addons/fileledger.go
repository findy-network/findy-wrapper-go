package addons

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

const fileName = "FINDY_FILE_LEDGER"

var filename = fullFilename(fileName)

// file is a ledger addon which implements transient ledger. It writes
// ledger data to JSON file and reads it from there. It's convenient for unit
// test and some development cases.
type file struct {
	Mem
}

func (m *file) Close() {
	// nothing to do in this version
}

func (m *file) Open(name ...string) bool {
	if strings.Contains(name[0], "=") {
		sub := strings.Split(name[0], "=")
		m.Mem.Open(sub[1])
		name[0] = sub[0]
	} else {
		m.Mem.Open("")
	}

	filename = fullFilename(name[0])
	glog.V(3).Infoln("-- file ledger:", filename)

	if fileExists() {
		try.To(m.load(filename))
	}
	return true
}

func (m *file) Write(tx plugin.TxInfo, ID, data string) (err error) {
	defer err2.Return(&err)

	try.To(m.Mem.Write(tx, ID, data))
	try.To(m.save(filename))

	return nil
}

func (m *file) Read(tx plugin.TxInfo, ID string) (name string, value string, err error) {
	return m.Mem.Read(tx, ID)
}

func (m *file) load(filename string) (err error) {
	defer err2.Return(&err)

	m.Mem.Mem.Lock()
	defer m.Mem.Mem.Unlock()

	if filename == "" {
		m.Mem.Mem.Ory = make(map[string]string)
		return nil
	}

	data := try.To1(os.ReadFile(filename))
	m.Mem.Mem.Ory = *newFromData(data)

	return nil
}

func (m *file) save(filename string) (err error) {
	defer err2.Return(&err)

	data := try.To1(json.MarshalIndent(m.Mem.Mem.Ory, "", "\t"))
	return writeJSONFile(filename, data)
}

func newFromData(data []byte) (r *map[string]string) {
	r = new(map[string]string)
	try.To(json.Unmarshal(data, r))
	return r
}

var fileLedger = &file{
	Mem: Mem{
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
	},
}

func init() {
	pool.RegisterPlugin(fileName, fileLedger)
}

func writeJSONFile(name string, json []byte) error {
	return os.WriteFile(name, json, 0644)
}

func fileExists() bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func fullFilename(fn ...string) string {
	const workerSubPath = "/.indy_client/"

	home := try.To1(user.Current()).HomeDir
	args := make([]string, len(fn)+2)
	args[0] = home
	args[1] = workerSubPath

	// first make sure we have proper base folder for our file
	pathBase := filepath.Join(args...)
	try.To(os.MkdirAll(pathBase, os.ModePerm)) // this panics if err

	// second build the whole file name by adding our filename args
	args = append(args, fn...)
	base := filepath.Join(args...)
	base += ".json"
	return base
}
