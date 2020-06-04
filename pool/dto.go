package pool

// Config is pool creating config structure used by indy.
type Config struct {
	// GenesisTxn is full filename of the genesis file.
	GenesisTxn string `json:"genesis_txn,omitempty"`
}
