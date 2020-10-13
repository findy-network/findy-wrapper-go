// Package file is a ledger addon which implements transient ledger. It writes
// ledger data to JSON file and reads it from there. It's convenient for unit
// test and some development cases.

package file

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

const addonName = "FINDY_FILE_LEDGER"

var filename = fullFilename(addonName)

type addon struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func (m *addon) Close() {
	//resetMemory()
}

func (m *addon) Open(name string) bool {
	filename = fullFilename(name)
	if fileExists() {
		err2.Check(m.load(filename))
	} else {
		resetMemory()
	}
	return true
}

func (m *addon) Write(ID, data string) (err error) {
	defer err2.Return(&err)
	m.mem.Lock()
	defer m.mem.Unlock()
	m.mem.ory[ID] = data
	err2.Check(m.save(filename))
	return nil
}

func (m *addon) Read(ID string) (name string, value string, err error) {
	m.mem.RLock()
	defer m.mem.RUnlock()
	return ID, m.mem.ory[ID], nil
}

func (m *addon) load(filename string) error {
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

func (m *addon) save(filename string) (err error) {
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

func writeJSONFile(name string, json []byte) error {
	err := ioutil.WriteFile(name, json, 0644)
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
