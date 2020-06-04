package c2go

//#include <stdio.h>
//#include <stdlib.h>
//#include "findy_glue.h"
import "C"
import (
	"unsafe"

	"github.com/optechlab/findy-go/internal/ctx"
)

func WalletGenerateKey(seed string) ctx.Channel {
	seedInC := C.CString(seed)
	defer C.free(unsafe.Pointer(seedInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("WalletGenerateKey")
	C.findy_generate_wallet_key(C.int(cmdHandle), seedInC)
	return ch
}

func WalletClose(handle int) ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.NamedPush("WalletClose")
	C.findy_close_wallet(C.int(cmdHandle), C.int(handle))
	return ch
}

func WalletCreate(config, credentials string) ctx.Channel {
	configC := C.CString(config)
	credentialsC := C.CString(credentials)
	defer C.free(unsafe.Pointer(configC))
	defer C.free(unsafe.Pointer(credentialsC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("CreateWallet")
	C.findy_create_wallet(C.int(cmdHandle), configC, credentialsC)
	return ch
}

func FindyDeleteWallet(config, credentials string) ctx.Channel {
	configC := C.CString(config)
	credentialsC := C.CString(credentials)
	defer C.free(unsafe.Pointer(configC))
	defer C.free(unsafe.Pointer(credentialsC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("DeleteWallet")
	C.findy_delete_wallet(C.int(cmdHandle), configC, credentialsC)
	return ch
}

func WalletOpen(config, credentials string) ctx.Channel {
	configC := C.CString(config)
	credentialsC := C.CString(credentials)
	defer C.free(unsafe.Pointer(configC))
	defer C.free(unsafe.Pointer(credentialsC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("WalletOpen")
	C.findy_open_wallet(C.int(cmdHandle), configC, credentialsC)
	return ch
}

func FindyImportWallet(config, credentials, importCfg string) ctx.Channel {
	configInC := C.CString(config)
	defer C.free(unsafe.Pointer(configInC))
	credentialsInC := C.CString(credentials)
	defer C.free(unsafe.Pointer(credentialsInC))
	importCfgInC := C.CString(importCfg)
	defer C.free(unsafe.Pointer(importCfgInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyImportWallet")
	C.findy_import_wallet(C.int(cmdHandle), configInC, credentialsInC, importCfgInC)
	return ch
}

func FindyExportWallet(wallet int, exportCfg string) ctx.Channel {
	exportCfgInC := C.CString(exportCfg)
	defer C.free(unsafe.Pointer(exportCfgInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyExportWallet")
	C.findy_export_wallet(C.int(cmdHandle), C.int(wallet), exportCfgInC)
	return ch
}
