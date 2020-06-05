package wallet

// StorageConfig is indy specific JSON configuration struct to set wallet
// location path.
type StorageConfig struct {
	Path string `json:"path,omitempty"`
}

// Config is a config struct for indy wallet functions.
type Config struct {
	// ID is the name of the wallet i.e. the last directory in the wallet path
	ID string `json:"id"`

	// StorageType is optional, only use if indy_register_wallet_storage() is
	// called. 'Default' value is for local files sytem wallets.
	StorageType string `json:"storage_type,omitempty"`

	// StorageConfig is optional, use when the wallet root path needs to be set.
	*StorageConfig `json:"storage_config,omitempty"`
}

// Credentials is a indy specific struct to set wallet credentials for the indy
// wallet functions for opening and exporting the wallet. The Path is only for
// export/import. The StorageConfig in the Config struct is for the actual
// wallet location.
type Credentials struct {
	// Path is used only with export and import functions
	Path string `json:"path,omitempty"`

	// Key with method are the actual credentials.
	Key string `json:"key,omitempty"`

	// KeyDerivationMethod is for the algorithm to use for the Key:
	//  ARGON2I_MOD - derive secured wallet master rekey (used by default)
	//  ARGON2I_INT - derive secured wallet master rekey (less secured but faster)
	//  RAW - raw wallet key master provided (skip derivation)
	//        RAW keys can be generated with indy_generate_wallet_key call
	KeyDerivationMethod string `json:"key_derivation_method,omitempty"`
}
