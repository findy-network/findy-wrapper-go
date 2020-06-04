package c2go

//#include <stdio.h>
//#include <stdlib.h>
//#include "findy_glue.h"
import "C"
import (
	"unsafe"

	"github.com/optechlab/findy-go/internal/ctx"
)

func CreateAndStoreMyDid(wallet int, didJSON string) ctx.Channel {
	didJSONInC := C.CString(didJSON)
	defer C.free(unsafe.Pointer(didJSONInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("CreateAndStoreMyDid")
	C.findy_create_and_store_my_did(C.int(cmdHandle), C.int(wallet), didJSONInC)
	return ch
}

func ListDids(wallet int) ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_list_my_dids_with_meta(C.int(cmdHandle), C.int(wallet))
	return ch
}

func KeyForDid(pool, wallet int, did string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_key_for_did(C.int(cmdHandle), C.int(pool), C.int(wallet), didInC)
	return ch
}

func FindyStoreTheirDid(wallet int, idJSON string) ctx.Channel {
	idJSONInC := C.CString(idJSON)
	defer C.free(unsafe.Pointer(idJSONInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyStoreTheirDid: " + idJSON)
	C.findy_store_their_did(C.int(cmdHandle), C.int(wallet), idJSONInC)
	return ch
}

func FindyKeyForLocalDid(wallet int, did string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyKeyForLocalDid: " + did)
	C.findy_key_for_local_did(C.int(cmdHandle), C.int(wallet), didInC)
	return ch
}

func FindySetDidMetadata(wallet int, did, meta string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	metaInC := C.CString(meta)
	defer C.free(unsafe.Pointer(metaInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_set_did_metadata(C.int(cmdHandle), C.int(wallet), didInC, metaInC)
	return ch
}

func FindyGetDidMetadata(wallet int, did string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_get_did_metadata(C.int(cmdHandle), C.int(wallet), didInC)
	return ch
}

func FindyGetMyDidWithMeta(wallet int, did string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_get_my_did_with_meta(C.int(cmdHandle), C.int(wallet), didInC)
	return ch
}

func FindySetEndpointForDid(wallet int, did, address, key string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	addressInC := C.CString(address)
	defer C.free(unsafe.Pointer(addressInC))
	keyInC := C.CString(key)
	defer C.free(unsafe.Pointer(keyInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("DidSetEndpoint: " + did)
	C.findy_set_endpoint_for_did(C.int(cmdHandle), C.int(wallet), didInC, addressInC, keyInC)
	return ch
}

func FindyGetEndpointForDid(wallet, pool int, did string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("DidGetEndpoint: " + did)
	C.findy_get_endpoint_for_did(C.int(cmdHandle), C.int(wallet), C.int(pool), didInC)
	return ch
}
