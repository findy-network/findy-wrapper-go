package did

// Did is a helper struct to build JSON data for indy functions.
type Did struct {
	Did        string `json:"did,omitempty"`
	Seed       string `json:"seed,omitempty"`
	CryptoType string `json:"crypto_type,omitempty"`
	Cid        bool   `json:"cid,omitempty"`
	VerKey     string `json:"verkey,omitempty"`
}
