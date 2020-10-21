package addons

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/lainio/err2"
)

const fileName = "FINDY_FILE_LEDGER"

var filename = fullFilename(fileName)

// file is a ledger addon which implements transient ledger. It writes
// ledger data to JSON file and reads it from there. It's convenient for unit
// test and some development cases.
type file struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func (m *file) Close() {
	//resetMemory()
}

func (m *file) Open(name string) bool {
	filename = fullFilename(name)
	if fileExists() {
		err2.Check(m.load(filename))
	} else {
		resetMemory()
	}
	return true
}

func (m *file) Write(ID, data string) (err error) {
	defer err2.Return(&err)
	m.mem.Lock()
	defer m.mem.Unlock()
	m.mem.ory[ID] = data
	err2.Check(m.save(filename))
	return nil
}

func (m *file) Read(ID string) (name string, value string, err error) {
	m.mem.RLock()
	defer m.mem.RUnlock()
	return ID, m.mem.ory[ID], nil
}

func (m *file) load(filename string) error {
	m.mem.Lock()
	defer m.mem.Unlock()
	if filename == "" {
		m.mem.ory = make(map[string]string)
		return nil
	}
	data, err := readJSONFile(filename)
	if err != nil {
		return err
	}
	m.mem.ory = *newFromData(data)
	return nil
}

func (m *file) save(filename string) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(m.mem.ory, "", "\t"); err != nil {
		return err
	}
	return writeJSONFile(filename, data)
}

func newFromData(data []byte) (r *map[string]string) {
	r = new(map[string]string)
	err := json.Unmarshal(data, r)
	if err != nil {
		panic(err)
	}
	return
}

var fileLedger = &file{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

func init() {
	pool.RegisterPlugin(fileName, fileLedger)
}

func resetMemory() {
	fileLedger.mem.ory = make(map[string]string)
}

func writeJSONFile(name string, b []byte) error {
	err := ioutil.WriteFile(name, b, 0644)
	return err
}

func readJSONFile(name string) ([]byte, error) {
	result, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return result, err
}

func fileExists() bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func fullFilename(fn ...string) string {
	const workerSubPath = "/.indy_client/"

	home := homeDir()
	args := make([]string, len(fn)+2)
	args[0] = home
	args[1] = workerSubPath

	// first make sure we have proper base folder for our file
	pathBase := filepath.Join(args...)
	err2.Check(os.MkdirAll(pathBase, os.ModePerm)) // this panics if err

	// second build the whole file name by adding our filename args
	args = append(args, fn...)
	base := filepath.Join(args...)
	base += ".json"
	return base
}

func homeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		err2.Check(err)
	}
	return currentUser.HomeDir
}
