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
	_ int,
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
	_,
	wallet int,
	submitter,
	credDef string,
) (err error) {
	defer err2.Handle(&err)

	return writePluginLedgers(
		plugin.TxInfo{
			TxType:       plugin.TxTypeCredDef,
			Wallet:       wallet,
			SubmitterDID: submitter,
		},
		credDef)
}

func writePluginLedgers(tx plugin.TxInfo, data string) error {
	var raw map[string]interface{}
	dto.FromJSONStr(data, &raw)
	dataID := raw["id"].(string)
	return pool.Write(tx, dataID, data)
}

// ReadSchema reads schema from ledgers by ID. If multiple ledger plugins is
// used, it returns where it can find data first.
func ReadSchema(_ int, submitter, ID string) (sID, s string, err error) {
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
	_ int,
	wallet int,
	submitter string,
	scJSON string,
) (err error) {
	defer err2.Handle(&err)

	return writePluginLedgers(
		plugin.TxInfo{
			TxType:       plugin.TxTypeSchema,
			Wallet:       wallet,
			SubmitterDID: submitter,
		},
		scJSON)
}

// WriteDID writes DID to ledger. If multiple ledger plugins in in use, it
// writes to all of them. Note! There are no ReadDID yet, because it is not used
// explicitly. Some of the indy SDK functions read ledger implicitly like
// did_get_key().
func WriteDID(
	_,
	wallet int,
	submitterDID,
	targetDID,
	verKey,
	alias,
	role string,
) (err error) {
	defer err2.Handle(&err)

	return pool.Write(
		plugin.TxInfo{
			TxType:       plugin.TxTypeDID,
			Wallet:       wallet,
			SubmitterDID: submitterDID,
			VerKey:       verKey,
			Alias:        alias,
			Role:         role,
		},
		submitterDID, targetDID)
}
