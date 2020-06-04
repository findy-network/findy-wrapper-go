// Package plugin is an interface package for ledger addons.
package plugin

// Plugin is a plugin interface for addon ledger implementations.
type Plugin interface {
	Open(name string) bool
	Close()
}

// Mapper is an property getter/setter interface for addon ledger
// implementations.
type Mapper interface {
	Write(ID, data string) error
	Read(ID string) (string, string, error)
}

// Ledger is a plugin interface used to offer implementations of addon ledgers.
// See pool package for more information.
type Ledger interface {
	Plugin
	Mapper
}
