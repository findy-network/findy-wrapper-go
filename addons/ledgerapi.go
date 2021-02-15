package addons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"

	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/lainio/err2"
)

const apiLedgerName = "FINDY_API_LEDGER"
const AuthToken = "aff62524-e503-11e9-81b4-2a2ae2dbcce4" // this must come from env
const BaseAddress = "http://localhost:3000/dev/"         // this must come from env
const DIDLengthMin = 21
const DIDLengthMax = 22

// api is a ledger addon which implements reading / writing data to the DB based ledger API.
// It writes ledger data to memory and reads it from there. It's convenient for unit test
// and some development cases.
type api struct {
	mem struct {
		sync.RWMutex
		ory map[string]string
	}
}

func CleanString(data string) string {
	// clean
	re := regexp.MustCompile(`\n`)
	data = re.ReplaceAllString(data, "")
	re = regexp.MustCompile(`\t`)
	data = re.ReplaceAllString(data, "")
	return data
}

func (a *api) Close() {
	// resetApiLedger()
}

func (a *api) Open(name string) bool {
	resetApiLedger()
	return name == apiLedgerName
}

func (a *api) Write(ID, data string) error {
	a.mem.Lock()
	defer a.mem.Unlock()

	const path = "store"

	// fmt.Println("dataa: ", data)
	// clean
	data = CleanString(data)
	// fmt.Println("dataa siivottuna: ", data)
	jsonEncodedData, err := json.Marshal(data)
	if err != nil {
		err2.Check(err)
	}
	// fmt.Println("dataa json encoodattuna: ", string(jsonEncodedData))

	// Build the http POST request towards ledger API
	client := &http.Client{}
	// address := fmt.Sprint(BaseAddress, path, ID)
	address := fmt.Sprint(BaseAddress, path)
	// fmt.Println("osoite: ", address)
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonEncodedData))

	// Set headers
	bearer := fmt.Sprint("Bearer ", AuthToken)
	req.Header.Add("authorization", bearer)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		err2.Check(err)
	}

	fmt.Println(resp.Status)
	a.mem.ory[ID] = data
	return nil
}

func (a *api) Read(ID string) (name string, value string, err error) {
	a.mem.RLock()
	defer a.mem.RUnlock()

	// expr := fmt.Sprintf("^[1-9A-HJ-NP-Za-km-z]{%d,%d}", DIDLengthMin, DIDLengthMax)
	var path string
	expr := fmt.Sprint(`:`)
	r, _ := regexp.Compile(expr)
	if r.MatchString(ID) {
		fmt.Println("ID", ID, "contains :")
		expr2 := fmt.Sprintf("^([^:]*:){3}[^:]*$") // check for 3 colons
		r, _ := regexp.Compile(expr2)
		if r.MatchString(ID) {
			// SCHEMA does have 3 colons
			fmt.Println("ID", ID, "contains 3 :. It is a SCHEMA")
			path = "schema/"
		} else {
			// CRED_DEF does have 7 colons
			fmt.Println("ID", ID, "contains more than 0 but not 3 :. It must be CREF_DEF")
			path = "cred_def/"
		}

	} else {
		// NYM txnId does not have any colons
		fmt.Println("ID", ID, "does not contain :. It's a NYM")
		path = "nym/"
	}

	// Build the http GET request towards ledger API
	client := &http.Client{}
	tempId := "FLfpAM9gWWLBZ4Lq6iGE9F:4:FLfpAM9gWWLBZ4Lq6iGE9F:3:CL:56310:3c939071-bd9d-4059-91fb-f7206af2fbce:CL_ACCUM:1-1024"
	// address := fmt.Sprint(BaseAddress, path, ID)
	address := fmt.Sprint(BaseAddress, path, tempId)
	req, err := http.NewRequest("GET", address, nil)

	// Set headers
	bearer := fmt.Sprint("Bearer ", AuthToken)
	req.Header.Add("authorization", bearer)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	// ToDO error handling

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

	return ID, a.mem.ory[ID], nil
}

var apiLedger = &api{mem: struct {
	sync.RWMutex
	ory map[string]string
}{}}

func init() {
	pool.RegisterPlugin(apiLedgerName, apiLedger)
}

func resetApiLedger() {
	apiLedger.mem.ory = make(map[string]string)
}
