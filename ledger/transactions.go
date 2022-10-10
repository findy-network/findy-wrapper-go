package ledger

import (
	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/lainio/err2"
)

// ReadCredDef reads cred def from ledgers by cred def ID. If multiple ledger
// plugins is used, it returns where it can find data first.
func ReadCredDef(
	handle int,
	submitter,
	credDefID string,
) (
	cdID,
	cd string,
	err error,
) {
	return pool.Read(
		plugin.TxInfo{
			TxType:       plugin.TxTypeCredDef,
			SubmitterDID: submitter,
		},
		credDefID)
}

// WriteCredDef writes cred def to ledger. If multiple ledger plugins is in use,
// it writes data to all of them.
func WriteCredDef(
	handle,
	wallet int,
	submitter,
	credDef string,
) (err error) {
	defer err2.Returnf(&err, "write cred def")

	writePluginLedgers(
		plugin.TxInfo{
			TxType:       plugin.TxTypeCredDef,
			Wallet:       wallet,
			SubmitterDID: submitter,
		},
		credDef)
	return nil
}

func writePluginLedgers(tx plugin.TxInfo, data string) {
	var raw map[string]interface{}
	dto.FromJSONStr(data, &raw)
	dataID := raw["id"].(string)
	pool.Write(tx, dataID, data)
}

// ReadSchema reads schema from ledgers by ID. If multiple ledger plugins is
// used, it returns where it can find data first.
func ReadSchema(handle int, submitter, ID string) (sID, s string, err error) {
	return pool.Read(
		plugin.TxInfo{
			TxType:       plugin.TxTypeSchema,
			SubmitterDID: submitter,
		},
		ID,
	)
}

// WriteSchema writes schema to ledger. If multiple ledger plugins is in use, it
// writes to all of them.
func WriteSchema(
	handle int,
	wallet int,
	submitter string,
	scJSON string,
) (err error) {
	defer err2.Returnf(&err, "write schema")

	writePluginLedgers(
		plugin.TxInfo{
			TxType:       plugin.TxTypeSchema,
			Wallet:       wallet,
			SubmitterDID: submitter,
		},
		scJSON)
	return nil
}

// WriteDID writes DID to ledger. If multiple ledger plugins in in use, it
// writes to all of them. Note! There are no ReadDID yet, because it is not used
// explicitly. Some of the indy SDK functions read ledger implicitly like
// did_get_key().
func WriteDID(
	handle,
	wallet int,
	submitterDID,
	targetDID,
	verKey,
	alias,
	role string,
) (err error) {
	defer err2.Returnf(&err, "write DID")

	pool.Write(
		plugin.TxInfo{
			TxType:       plugin.TxTypeDID,
			Wallet:       wallet,
			SubmitterDID: submitterDID,
			VerKey:       verKey,
			Alias:        alias,
			Role:         role,
		},
		submitterDID, targetDID)

	return nil
}
