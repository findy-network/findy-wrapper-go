package wallet

import (
	"fmt"

	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/findy-network/findy-wrapper-go/internal/c2go"
	"github.com/findy-network/findy-wrapper-go/internal/ctx"
)

// GenerateKey generates a wallet master key. Returned key is compatible with
// "RAW" key derivation method. It allows to avoid expensive key derivation for
// use cases when wallet keys can be stored in a secure enclave.
func GenerateKey(seed string) ctx.Channel {
	if seed == "" {
		seed = "{}"
	} else {
		seed = fmt.Sprintf("{\"seed\":\"%s\"}", seed)
	}
	return c2go.WalletGenerateKey(seed)
}

// Close closes the wallet identified by wallet handle.
func Close(handle int) ctx.Channel {
	return c2go.WalletClose(handle)
}

// Create creates the wallet. See Config and Credentials morey info.
func Create(config Config, credentials Credentials) ctx.Channel {
	return c2go.WalletCreate(dto.ToJSON(config), dto.ToJSON(credentials))
}

// Delete permanently deletes the wallet with all of its file structures.
func Delete(config Config, credentials Credentials) ctx.Channel {
	return c2go.FindyDeleteWallet(dto.ToJSON(config), dto.ToJSON(credentials))
}

// Open opens a wallet identified by Config and Credentials. See more info of
// them. The successful wallet handle is Result.Handle().
func Open(config Config, credentials Credentials) ctx.Channel {
	return c2go.WalletOpen(dto.ToJSON(config), dto.ToJSON(credentials))
}

// Export generates a single encrypted file of the wallet identified by Config
// and Credentials. See more info of them.
func Export(handle int, exportCfg Credentials) ctx.Channel {
	return c2go.FindyExportWallet(handle, dto.ToJSON(exportCfg))
}

// Import creates a new wallet from previously exported (Export function)
// wallet's export file. See more info of Config and Credentials.
func Import(config Config, credentials, importCfg Credentials) ctx.Channel {
	return c2go.FindyImportWallet(dto.ToJSON(config), dto.ToJSON(credentials), dto.ToJSON(importCfg))
}
