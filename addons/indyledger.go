package addons

import (
	"errors"

	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/findy-network/findy-wrapper-go/internal/c2go"
	"github.com/findy-network/findy-wrapper-go/ledger"
	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
)

const indyLedgerAddonName = "FINDY_LEDGER"

// Indy is a ledger addon which implements real Indy ledger pool client.
// The machine which uses it must have indy pool named FINDY_LEDGER.
// Also it's important that memory or file ledger is used same time because both
// writing and reading to addon is done asynchronously.
type Indy struct {
	handle int
}

func (ao *Indy) Close() {
	c2go.PoolCloseLedger(ao.handle)
}

func (ao *Indy) Open(name ...string) (ok bool) {
	poolName := name[0]
	if poolName == "" {
		poolName = indyLedgerAddonName
	}

	defer err2.Catch(func(err error) {
		ok = false
		glog.Errorln("cannot open", indyLedgerAddonName, "by name", poolName)
		glog.Errorln(err)
	})

	r := <-c2go.PoolOpenLedger(poolName)
	try.To(r.Err())
	ao.handle = r.Handle()
	return true
}

func (ao *Indy) Write(tx plugin.TxInfo, ID, data string) error {
	switch tx.TxType {
	case plugin.TxTypeDID:
		return ao.WriteDID(tx, ID, data)

	case plugin.TxTypeSchema:
		return ao.WriteSchema(tx, ID, data)

	case plugin.TxTypeCredDef:
		return ao.WriteCredDef(tx, ID, data)

	}

	return nil
}

func (ao *Indy) Read(tx plugin.TxInfo, ID string) (name string, value string, err error) {
	switch tx.TxType {
	case plugin.TxTypeDID:
		assert.D.True(false, "we don't support DID reading from ledger")

	case plugin.TxTypeSchema:
		return ao.ReadSchema(tx, ID)

	case plugin.TxTypeCredDef:
		return ao.ReadCredDef(tx, ID)
	}
	return
}

var indyAddonLedger = new(Indy)

func init() {
	pool.RegisterPlugin(indyLedgerAddonName, indyAddonLedger)
}

func (ao *Indy) ReadCredDef(
	tx plugin.TxInfo,
	credDefID string,
) (name string, value string, err error) {
	defer err2.Return(&err)

	glog.V(100).Infoln("submitter:", tx.SubmitterDID)

	r := <-ledger.BuildGetCredDefRequest(tx.SubmitterDID, credDefID)
	try.To(r.Err())

	r = <-ledger.SubmitRequest(ao.handle, r.Str1())
	try.To(r.Err())

	lresp := r.Str1()

	// parse ledger response to credDefID and credDef
	r = <-ledger.ParseGetCredDefResponse(lresp)
	try.To(r.Err())

	name = r.Str1()
	assert.P.True(name == credDefID)

	credDef := r.Str2()
	return credDefID, credDef, nil
}

func (ao *Indy) ReadSchema(
	tx plugin.TxInfo,
	ID string,
) (name string, value string, err error) {
	defer err2.Return(&err)

	glog.V(100).Infoln("submitter:", tx.SubmitterDID)

	r := <-ledger.BuildGetSchemaRequest(tx.SubmitterDID, ID)
	try.To(r.Err())

	gsr := r.Str1()
	r = <-ledger.SubmitRequest(ao.handle, gsr)
	try.To(r.Err())

	gsr = r.Str1()
	r = <-ledger.ParseGetSchemaResponse(gsr)
	try.To(r.Err())

	readSchemaID := r.Str1()
	scJSON := r.Str2()

	return readSchemaID, scJSON, nil
}

func (ao *Indy) WriteDID(
	tx plugin.TxInfo,
	ID string,
	data string,
) (err error) {
	defer err2.Return(&err)

	glog.V(1).Infoln("submitter:", tx.SubmitterDID)

	r := <-ledger.BuildNymRequest(tx.SubmitterDID, data, tx.VerKey, tx.Alias, tx.Role)
	try.To(r.Err())

	req := r.Str1()
	r = <-ledger.SignAndSubmitRequest(ao.handle, tx.Wallet, tx.SubmitterDID, req)
	try.To(r.Err())
	try.To(checkWriteResponse(r.Str1()))

	return nil
}

func (ao *Indy) WriteSchema(
	tx plugin.TxInfo,
	ID string,
	data string,
) (err error) {
	defer err2.Return(&err)

	glog.V(1).Infoln("submitter:", tx.SubmitterDID)

	r := <-ledger.BuildSchemaRequest(tx.SubmitterDID, data)
	try.To(r.Err())

	srq := r.Str1()
	r = <-ledger.SignAndSubmitRequest(ao.handle, tx.Wallet, tx.SubmitterDID, srq)
	try.To(r.Err())

	try.To(checkWriteResponse(r.Str1()))
	return nil
}

func (ao *Indy) WriteCredDef(
	tx plugin.TxInfo,
	ID string,
	data string,
) (err error) {
	defer err2.Return(&err)

	glog.V(1).Infoln("submitter:", tx.SubmitterDID)

	r := <-ledger.BuildCredDefRequest(tx.SubmitterDID, data)
	try.To(r.Err())

	r = <-ledger.SignAndSubmitRequest(ao.handle, tx.Wallet, tx.SubmitterDID, r.Str1())
	try.To(r.Err())
	try.To(checkWriteResponse(r.Str1()))
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
