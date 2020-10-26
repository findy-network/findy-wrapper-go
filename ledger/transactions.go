package ledger

import (
	"errors"

	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

// ReadCredDef reads cred def from ledgers by cred def ID. If multiple ledger
// plugins is used, it returns where it can find data first.
func ReadCredDef(handle int, submitter,
	credDefID string) (cdID, cd string, err error) {

	defer err2.Annotate("read cred def", &err)

	cdID, cd, err = pool.Read(credDefID)
	if !pool.IsIndyLedgerOpen(handle) {
		return
	}

	// if we still need to read from the indy ledger. In future we could start
	// to fill other ledger plugins to make poor mans cache
	if cd == "" {
		r := <-BuildGetCredDefRequest(submitter, credDefID)
		err2.Check(r.Err())

		r = <-SubmitRequest(handle, r.Str1())
		err2.Check(r.Err())

		lresp := r.Str1()

		// parse ledger response to credDefID and credDef
		r = <-ParseGetCredDefResponse(lresp)
		err2.Check(r.Err())

		cdID = r.Str1()
		if credDefID != cdID { // Just debugging and reverse-engineering
			glog.Info("=========== CRED DEF ID IS NOT SAME")
		}
		credDef := r.Str2()
		return credDefID, credDef, nil
	}
	return "", "", nil
}

// WriteCredDef writes cred def to ledger. If multiple ledger plugins is in use,
// it writes data to all of them.
func WriteCredDef(handle, wallet int, submitter, credDef string) (err error) {
	defer err2.Annotate("write cred def", &err)

	writePluginLedgers(credDef)
	if !pool.IsIndyLedgerOpen(handle) {
		return nil
	}

	r := <-BuildCredDefRequest(submitter, credDef)
	err2.Check(r.Err())

	r = <-SignAndSubmitRequest(handle, wallet, submitter, r.Str1())
	err2.Check(r.Err())
	err2.Check(checkWriteResponse(r.Str1()))

	return nil
}

func writePluginLedgers(data string) {
	var raw map[string]interface{}
	dto.FromJSONStr(data, &raw)
	credDefID := raw["id"].(string)
	pool.Write(credDefID, data)
}

// ReadSchema reads schema from ledgers by ID. If multiple ledger plugins is
// used, it returns where it can find data first.
func ReadSchema(handle int, submitter, ID string) (sID, s string, err error) {
	defer err2.Annotate("read schema", &err)

	sID, s, err = pool.Read(ID)
	if !pool.IsIndyLedgerOpen(handle) {
		return
	}

	// if we still need to read from the indy ledger. In future we could start
	// to fill other ledger plugins to make poor mans cache
	if s == "" {
		r := <-BuildGetSchemaRequest(submitter, ID)
		err2.Check(r.Err())

		gsr := r.Str1()
		r = <-SubmitRequest(handle, gsr)
		err2.Check(r.Err())

		gsr = r.Str1()
		r = <-ParseGetSchemaResponse(gsr)
		err2.Check(r.Err())

		readSchemaID := r.Str1()
		scJSON := r.Str2()

		return readSchemaID, scJSON, nil
	}
	return "", "", nil
}

// WriteSchema writes schema to ledger. If multiple ledger plugins is in use, it
// writes to all of them.
func WriteSchema(handle int, wallet int, submitter string, scJSON string) (err error) {
	defer err2.Annotate("write schema", &err)

	writePluginLedgers(scJSON)
	if !pool.IsIndyLedgerOpen(handle) {
		return nil
	}

	r := <-BuildSchemaRequest(submitter, scJSON)
	err2.Check(r.Err())

	srq := r.Str1()
	r = <-SignAndSubmitRequest(handle, wallet, submitter, srq)
	err2.Check(r.Err())

	err2.Check(checkWriteResponse(r.Str1()))
	return nil
}

func checkWriteResponse(r string) error {
	type response struct {
		Op         string `json:"op"`
		Identifier string `json:"identifier"`
		ReqID      uint64 `json:"reqId"`
		Reason     string `json:"reason"`
	}
	var res response
	dto.FromJSONStr(r, &res)

	switch res.Op {
	case "REJECT":
		return errors.New(res.Reason)
	case "REPLY": // we know this one, it's here for debugging
		return nil
	default:
		return nil
	}
}

// WriteDID writes DID to ledger. If multiple ledger plugins in in use, it
// writes to all of them. Note! There are no ReadDID yet, because it is not used
// explicitly. Some of the indy SDK functions read ledger implicitly like
// did_get_key().
func WriteDID(handle, wallet int,
	submitterDID, targetDID, verKey, alias, role string) (err error) {

	defer err2.Annotate("write DID", &err)

	pool.Write(submitterDID, targetDID)
	if !pool.IsIndyLedgerOpen(handle) {
		return nil
	}

	r := <-BuildNymRequest(submitterDID, targetDID, verKey, alias, role)
	err2.Check(r.Err())

	req := r.Str1()
	r = <-SignAndSubmitRequest(handle, wallet, submitterDID, req)
	err2.Check(r.Err())
	err2.Check(checkWriteResponse(r.Str1()))

	return nil
}
