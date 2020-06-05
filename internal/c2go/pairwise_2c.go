package c2go

//#include <stdio.h>
//#include <stdlib.h>
//#include "findy_glue.h"
import "C"
import (
	"unsafe"

	"github.com/findy-network/findy-wrapper-go/internal/ctx"
)

func FindyIsPairwiseExists(wallet int, theirDID string) ctx.Channel {
	theirDIDInC := C.CString(theirDID)
	defer C.free(unsafe.Pointer(theirDIDInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_is_pairwise_exists(C.int(cmdHandle), C.int(wallet), theirDIDInC)
	return ch
}

func FindyCreatePairwise(wallet int, theirDID, myDID, metaData string) ctx.Channel {
	theirDIDInC := C.CString(theirDID)
	defer C.free(unsafe.Pointer(theirDIDInC))
	myDIDInC := C.CString(myDID)
	defer C.free(unsafe.Pointer(myDIDInC))
	metaDataInC := C.CString(metaData)
	defer C.free(unsafe.Pointer(metaDataInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_create_pairwise(C.int(cmdHandle), C.int(wallet), theirDIDInC, myDIDInC, metaDataInC)
	return ch
}

func FindyListPairwise(wallet int) ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_list_pairwise(C.int(cmdHandle), C.int(wallet))
	return ch
}

func FindyGetPairwise(wallet int, theirDID string) ctx.Channel {
	theirDIDInC := C.CString(theirDID)
	defer C.free(unsafe.Pointer(theirDIDInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_get_pairwise(C.int(cmdHandle), C.int(wallet), theirDIDInC)
	return ch
}

func FindySetPairwiseMetadata(wallet int, theirDID, metaData string) ctx.Channel {
	theirDIDInC := C.CString(theirDID)
	defer C.free(unsafe.Pointer(theirDIDInC))
	metaDataInC := C.CString(metaData)
	defer C.free(unsafe.Pointer(metaDataInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_set_pairwise_metadata(C.int(cmdHandle), C.int(wallet), theirDIDInC, metaDataInC)
	return ch
}
