// Package plugin is an interface package for ledger addons.
package plugin

// Plugin is a plugin interface for addon ledger implementations.
type Plugin interface {
	Open(name ...string) bool
	Close()
}

type TxType int

const (
	TxTypeDID TxType = iota
	TxTypeSchema
	TxTypeCredDef
)

type TxInfo struct {
	TxType

	Wallet       int
	SubmitterDID string
	VerKey       string
	Alias        string
	Role         string
}

var (
	TxDID     = TxInfo{TxType: TxTypeDID}
	TxSchema  = TxInfo{TxType: TxTypeSchema}
	TxCredDef = TxInfo{TxType: TxTypeCredDef}
)

// Mapper is an property getter/setter interface for addon ledger
// implementations.
type Mapper interface {
	Write(tx TxInfo, ID, data string) error
	Read(tx TxInfo, ID string) (string, string, error)
}

// Ledger is a plugin interface used to offer implementations of addon ledgers.
// See pool package for more information.
type Ledger interface {
	Plugin
	Mapper
}
