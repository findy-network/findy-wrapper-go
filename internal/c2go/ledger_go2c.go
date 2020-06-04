package c2go

//#include <stdio.h>
//#include <stdlib.h>
//#include "findy_glue.h"
import "C"
import (
	"unsafe"

	"github.com/optechlab/findy-go/internal/ctx"

	"github.com/optechlab/findy-go"
)

// MARK: ledger helper functions

func LedgerBuildNymRequest(submitterDid, targetDid, verkey, alias, role string) ctx.Channel {
	submitterDidInC := C.CString(submitterDid)
	targetDidInC := C.CString(targetDid)
	verkeyInC := C.CString(verkey)
	var aliasInC *C.char = C.findy_null_string
	if alias != findy.NullString {
		aliasInC = C.CString(alias)
		defer C.free(unsafe.Pointer(aliasInC))
	}
	var roleInC *C.char = C.findy_null_string
	if role != findy.NullString {
		roleInC = C.CString(role)
		defer C.free(unsafe.Pointer(roleInC))
	}
	defer C.free(unsafe.Pointer(submitterDidInC))
	defer C.free(unsafe.Pointer(targetDidInC))
	defer C.free(unsafe.Pointer(verkeyInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("LedgerBuildNymRequest")
	C.findy_build_nym_request(C.int(cmdHandle), submitterDidInC, targetDidInC, verkeyInC, aliasInC, roleInC)
	return ch
}

func LedgerSignAndSubmitRequest(pool, wallet int, submitterDid, request string) ctx.Channel {
	submitterDidInC := C.CString(submitterDid)
	requestInC := C.CString(request)
	defer C.free(unsafe.Pointer(submitterDidInC))
	defer C.free(unsafe.Pointer(requestInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("LedgerSignAndSubmitRequest")
	C.findy_sign_and_submit_request(C.int(cmdHandle), C.int(pool), C.int(wallet), submitterDidInC, requestInC)
	return ch
}

func LedgerBuildSchemaRequest(submitterDid, data string) ctx.Channel {
	submitterDidInC := C.CString(submitterDid)
	dataInC := C.CString(data)
	defer C.free(unsafe.Pointer(submitterDidInC))
	defer C.free(unsafe.Pointer(dataInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("LedgerBuildSchemaRequest")
	C.findy_build_schema_request(C.int(cmdHandle), submitterDidInC, dataInC)
	return ch
}

func FindyBuildGetNymRequest(submitterDid, targetDid string) ctx.Channel {
	submitterDidInC := C.CString(submitterDid)
	defer C.free(unsafe.Pointer(submitterDidInC))
	targetDidInC := C.CString(targetDid)
	defer C.free(unsafe.Pointer(targetDidInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyBuildGetNymRequest")
	C.findy_build_get_nym_request(C.int(cmdHandle), submitterDidInC, targetDidInC)
	return ch
}

func FindySubmitRequest(poolHandle int, requestJSON string) ctx.Channel {
	requestJSONInC := C.CString(requestJSON)
	defer C.free(unsafe.Pointer(requestJSONInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindySubmitRequest")
	C.findy_submit_request(C.int(cmdHandle), C.int(poolHandle), requestJSONInC)
	return ch
}

func FindyBuildAttribRequest(sDID, tDID, hasc, raw, enc string) ctx.Channel {
	var sDIDInC *C.char = C.findy_null_string
	if sDID != findy.NullString {
		sDIDInC = C.CString(sDID)
		defer C.free(unsafe.Pointer(sDIDInC))
	}
	tDIDInC := C.CString(tDID)
	defer C.free(unsafe.Pointer(tDIDInC))
	var hascInC *C.char = C.findy_null_string
	if hasc != findy.NullString {
		hascInC = C.CString(hasc)
		defer C.free(unsafe.Pointer(hascInC))
	}
	var rawInC *C.char = C.findy_null_string
	if raw != findy.NullString {
		rawInC = C.CString(raw)
		defer C.free(unsafe.Pointer(rawInC))
	}
	var encInC *C.char = C.findy_null_string
	if enc != findy.NullString {
		encInC = C.CString(enc)
		defer C.free(unsafe.Pointer(encInC))
	}
	cmdHandle, ch := ctx.CmdContext.NamedPush("AttrReq: " + raw)
	C.findy_build_attrib_request(C.int(cmdHandle), sDIDInC, tDIDInC, hascInC, rawInC, encInC)
	return ch
}

func FindyBuildGetAttribRequest(sDID, tDID, hasc, raw, enc string) ctx.Channel {
	var sDIDInC *C.char = C.findy_null_string
	if sDID != findy.NullString {
		sDIDInC = C.CString(sDID)
		defer C.free(unsafe.Pointer(sDIDInC))
	}
	tDIDInC := C.CString(tDID)
	defer C.free(unsafe.Pointer(tDIDInC))
	var hascInC *C.char = C.findy_null_string
	if hasc != findy.NullString {
		hascInC = C.CString(hasc)
		defer C.free(unsafe.Pointer(hascInC))
	}
	var rawInC *C.char = C.findy_null_string
	if raw != findy.NullString {
		rawInC = C.CString(raw)
		defer C.free(unsafe.Pointer(rawInC))
	}
	var encInC *C.char = C.findy_null_string
	if enc != findy.NullString {
		encInC = C.CString(enc)
		defer C.free(unsafe.Pointer(encInC))
	}
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyBuildGetAttribRequest")
	C.findy_build_get_attrib_request(C.int(cmdHandle), sDIDInC, tDIDInC, hascInC, rawInC, encInC)
	return ch
}

func FindyBuildCredDefRequest(submitter, data string) ctx.Channel {
	submitterInC := C.CString(submitter)
	defer C.free(unsafe.Pointer(submitterInC))
	dataInC := C.CString(data)
	defer C.free(unsafe.Pointer(dataInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyBuildCredDefRequest")
	C.findy_build_cred_def_request(C.int(cmdHandle), submitterInC, dataInC)
	return ch
}

func FindyBuildGetSchemaRequest(submitter, id string) ctx.Channel {
	submitterInC := C.CString(submitter)
	defer C.free(unsafe.Pointer(submitterInC))
	idInC := C.CString(id)
	defer C.free(unsafe.Pointer(idInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyBuildGetSchemaRequest")
	C.findy_build_get_schema_request(C.int(cmdHandle), submitterInC, idInC)
	return ch
}

func FindyParseGetSchemaResponse(response string) ctx.Channel {
	responseInC := C.CString(response)
	defer C.free(unsafe.Pointer(responseInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyParseGetSchemaResponse")
	C.findy_parse_get_schema_response(C.int(cmdHandle), responseInC)
	return ch
}

func FindyParseGetCredDefResponse(credDefResp string) ctx.Channel {
	credDefRespInC := C.CString(credDefResp)
	defer C.free(unsafe.Pointer(credDefRespInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyParseGetCredDefResponse")
	C.findy_parse_get_cred_def_response(C.int(cmdHandle), credDefRespInC)
	return ch
}

func FindyBuildGetCredDefRequest(submitter, id string) ctx.Channel {
	var submitterInC *C.char = C.findy_null_string
	if submitter != findy.NullString {
		submitterInC = C.CString(submitter)
		defer C.free(unsafe.Pointer(submitterInC))
	}
	idInC := C.CString(id)
	defer C.free(unsafe.Pointer(idInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyBuildGetCredDefRequest")
	C.findy_build_get_cred_def_request(C.int(cmdHandle), submitterInC, idInC)
	return ch
}
