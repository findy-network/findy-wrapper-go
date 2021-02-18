package addons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/lainio/err2"
)

const apiLedgerName = "FINDY_API_LEDGER"

// api is a ledger addon which implements reading / writing data to the DB based ledger API.
// It writes ledger data to memory before returning it and if it's stored in memory it serves
// the data from there instead of fetching it from the ledger API
type api struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

type nymTransaction struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

// clean possible line breaks and tabs from the data
func CleanString(data string) string {
	re := regexp.MustCompile(`\n`)
	data = re.ReplaceAllString(data, "")
	re = regexp.MustCompile(`\t`)
	data = re.ReplaceAllString(data, "")
	return data
}

// Check whether a string is a json or not
func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func (a *api) Close() {
	// resetAPILedger()
}

func (a *api) Open(name string) bool {
	resetAPILedger()
	return name == apiLedgerName
}

func (a *api) Write(ID, data string) (err error) {
	defer err2.Return(&err)
	a.mem.Lock()
	defer a.mem.Unlock()

	// Read APi Ledger Auth Token and base address from env
	authToken := os.Getenv("AuthToken")
	baseAddress := os.Getenv("BaseAddress")

	const path = "store"
	// fmt.Println("dataa: ", data)
	data = CleanString(data) // do some data cleaning if needed
	fmt.Println("dataa siivottuna: ", data)
	a.mem.ory[ID] = string(data) // store the data to the memory cache

	var jsonEncodedData []byte

	if isJSON(data) {
		// data contain json string so it must be either schema or cred def
		var err error
		jsonEncodedData, err = json.Marshal(data)
		if err != nil {
			err2.Check(err)
		}
		fmt.Println("dataa json encoodattuna: ", string(jsonEncodedData))
	} else {
		// data does not contain json string so it must be nym hash
		// put the data into json, which contains ID and dataType
		nymJson := &nymTransaction{
			Id:   ID,
			Data: data,
		}

		var err error
		jsonEncodedData, err = json.Marshal(nymJson)
		if err != nil {
			err2.Check(err)
		}
		// fmt.Println("dataa json encoodattuna: ", string(jsonEncodedData))
	}

	// Build the http POST request towards ledger API
	client := &http.Client{}
	address := fmt.Sprint(baseAddress, path)
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonEncodedData))
	if err != nil {
		err2.Check(err)
	}

	// Set headers
	bearer := fmt.Sprint("Bearer ", authToken)
	req.Header.Add("authorization", bearer)
	req.Header.Add("Accept", "application/json")

	// send the data
	resp, err := client.Do(req)
	if err != nil {
		err2.Check(err)
	}

	fmt.Println(resp.Status)
	return nil
}

func (a *api) Read(ID string) (name string, value string, err error) {
	// chekck if we have data in mem cache
	a.mem.RLock()
	if item, ok := a.mem.ory[ID]; ok {
		// data can be found from memcache, return it
		defer a.mem.RUnlock()
		return ID, item, nil
	}
	a.mem.RUnlock()

	// Read API Ledger Auth Token and base address from env
	authToken := os.Getenv("AuthToken")
	baseAddress := os.Getenv("BaseAddress")

	var path string
	expr := ":"
	r, _ := regexp.Compile(expr)
	if r.MatchString(ID) {
		expr2 := "^([^:]*:){3}[^:]*$" // check for 3 colon letters
		r, _ := regexp.Compile(expr2)
		if r.MatchString(ID) {
			// SCHEMA does have 3 colon letters
			// fmt.Println("ID", ID, "contains 3 :. It is a SCHEMA")
			path = "schema/"
		} else {
			// CRED_DEF does have 7 colon letters
			// fmt.Println("ID", ID, "contains more than 0 but not 3 :. It must be CREF_DEF")
			path = "cred_def/"
		}

	} else {
		// NYM txnId does not have any colon letters
		// fmt.Println("ID", ID, "does not contain :. It's a NYM")
		path = "nym/"
	}

	// Build the http GET request towards ledger API
	client := &http.Client{}
	address := fmt.Sprint(baseAddress, path, ID)
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		err2.Check(err)
	}

	// Set headers
	bearer := fmt.Sprint("Bearer ", authToken)
	req.Header.Add("authorization", bearer)
	req.Header.Add("Accept", "application/json")

	// handle request
	resp, err := client.Do(req)
	if err != nil {
		err2.Check(err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err2.Check(err)
	}

	fmt.Println(resp.Status)
	fmt.Println("response: ", string(respBody))
	a.mem.Lock()
	defer a.mem.Unlock()
	if path == "nym/" {
		// special handling. Pick just the hash from data-field.
		var nym *nymTransaction
		err := json.Unmarshal([]byte(respBody), &nym)
		if err != nil {
			err2.Check(err)
		}
		a.mem.ory[ID] = nym.Data // store the data to the memory cache
		fmt.Println("tallensin nyt sitten dataa ", a.mem.ory[ID])
		return ID, nym.Data, nil
	}
	a.mem.ory[ID] = string(respBody) // store the data to the memory cache
	fmt.Println("tallensin nyt sitten dataa ", a.mem.ory[ID])
	return ID, string(respBody), nil
}

var apiLedger = &api{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

func init() {
	// API Ledger token. Probably these will be set where this addon is called
	// So these two lines will be removed from here eventually
	os.Setenv("AuthToken", "aff62524-e503-11e9-81b4-2a2ae2dbcce4")
	os.Setenv("BaseAddress", "http://localhost:3000/dev/")
	pool.RegisterPlugin(apiLedgerName, apiLedger)
}

func resetAPILedger() {
	apiLedger.mem.ory = make(map[string]string)
}
