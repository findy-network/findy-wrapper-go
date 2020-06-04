package c2go

//#include <stdio.h>
//#include <stdlib.h>
//#include "findy_glue.h"
import "C"
import (
	"unsafe"

	"github.com/optechlab/findy-go"
	"github.com/optechlab/findy-go/internal/ctx"
)

func FindyIssuerCreateSchema(did, name, version, attrNames string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	nameInC := C.CString(name)
	defer C.free(unsafe.Pointer(nameInC))
	versionInC := C.CString(version)
	defer C.free(unsafe.Pointer(versionInC))
	attrNamesInC := C.CString(attrNames)
	defer C.free(unsafe.Pointer(attrNamesInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("CreateSchema: did:" + did)
	C.findy_issuer_create_schema(C.int(cmdHandle), didInC, nameInC, versionInC, attrNamesInC)
	return ch
}

func FindyIssuerCreateAndStoreCredentialDef(wallet int, did, schema, tag, sigType, config string) ctx.Channel {
	didInC := C.CString(did)
	defer C.free(unsafe.Pointer(didInC))
	schemaInC := C.CString(schema)
	defer C.free(unsafe.Pointer(schemaInC))
	tagInC := C.CString(tag)
	defer C.free(unsafe.Pointer(tagInC))
	var sigTypeInC *C.char = C.findy_null_string
	if sigType != findy.NullString {
		sigTypeInC = C.CString(sigType)
		defer C.free(unsafe.Pointer(sigTypeInC))
	}
	var configInC *C.char = C.findy_null_string
	if config != findy.NullString {
		configInC = C.CString(config)
		defer C.free(unsafe.Pointer(configInC))
	}
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyIssuerCreateAndStoreCredentialDef")
	C.findy_issuer_create_and_store_credential_def(C.int(cmdHandle), C.int(wallet), didInC, schemaInC, tagInC, sigTypeInC, configInC)
	return ch
}

func FindyIssuerCreateCredentialOffer(wallet int, credDefID string) ctx.Channel {
	credDefIDInC := C.CString(credDefID)
	defer C.free(unsafe.Pointer(credDefIDInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyIssuerCreateCredentialOffer")
	C.findy_issuer_create_credential_offer(C.int(cmdHandle), C.int(wallet), credDefIDInC)
	return ch
}

func FindyIssuerCreateCredential(wallet int, credOffer, credReq, credValues, revRegID string, blobHandle int) ctx.Channel {
	credOfferInC := C.CString(credOffer)
	defer C.free(unsafe.Pointer(credOfferInC))
	credReqInC := C.CString(credReq)
	defer C.free(unsafe.Pointer(credReqInC))
	credValuesInC := C.CString(credValues)
	defer C.free(unsafe.Pointer(credValuesInC))
	var revRegIDInC *C.char = C.findy_null_string
	if revRegID != findy.NullString {
		revRegIDInC = C.CString(revRegID)
		defer C.free(unsafe.Pointer(revRegIDInC))
	}
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyIssuerCreateCredential")
	C.findy_issuer_create_credential(C.int(cmdHandle), C.int(wallet), credOfferInC, credReqInC, credValuesInC, revRegIDInC, C.int(blobHandle))
	return ch
}

func FindyProverCreateMasterSecret(wallet int, id string) ctx.Channel {
	var idInC *C.char = C.findy_null_string
	if id != findy.NullString {
		idInC = C.CString(id)
		defer func() {
			C.free(unsafe.Pointer(idInC))
		}()
	}
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyProverCreateMasterSecret")
	C.findy_prover_create_master_secret(C.int(cmdHandle), C.int(wallet), idInC)
	return ch
}

func FindyProverCreateProof(wallet int, proofReq, reqCred, masterSec, schemas, credDef, revStates string) ctx.Channel {
	proofReqInC := C.CString(proofReq)
	defer C.free(unsafe.Pointer(proofReqInC))
	reqCredInC := C.CString(reqCred)
	defer C.free(unsafe.Pointer(reqCredInC))
	masterSecInC := C.CString(masterSec)
	defer C.free(unsafe.Pointer(masterSecInC))
	schemasInC := C.CString(schemas)
	defer C.free(unsafe.Pointer(schemasInC))
	credDefInC := C.CString(credDef)
	defer C.free(unsafe.Pointer(credDefInC))
	revStatesInC := C.CString(revStates)
	defer C.free(unsafe.Pointer(revStatesInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyProverCreateProof")
	C.findy_prover_create_proof(C.int(cmdHandle), C.int(wallet), proofReqInC, reqCredInC, masterSecInC, schemasInC, credDefInC, revStatesInC)
	return ch
}

func FindyProverCreateCredentialReq(wallet int, prover, credOffer, credDef, master string) ctx.Channel {
	proverInC := C.CString(prover)
	defer C.free(unsafe.Pointer(proverInC))
	credOfferInC := C.CString(credOffer)
	defer C.free(unsafe.Pointer(credOfferInC))
	credDefInC := C.CString(credDef)
	defer C.free(unsafe.Pointer(credDefInC))
	masterInC := C.CString(master)
	defer C.free(unsafe.Pointer(masterInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyProverCreateCredentialReq")
	C.findy_prover_create_credential_req(C.int(cmdHandle), C.int(wallet), proverInC, credOfferInC, credDefInC, masterInC)
	return ch
}

func FindyProverStoreCredential(wallet int, credID, credReqMeta, credentials, credDef, revRegDef string) ctx.Channel {
	var credIDInC *C.char = C.findy_null_string
	if credID != findy.NullString {
		credIDInC = C.CString(credID)
		defer C.free(unsafe.Pointer(credIDInC))
	}
	credReqMetaInC := C.CString(credReqMeta)
	defer C.free(unsafe.Pointer(credReqMetaInC))
	credentialsInC := C.CString(credentials)
	defer C.free(unsafe.Pointer(credentialsInC))
	credDefInC := C.CString(credDef)
	defer C.free(unsafe.Pointer(credDefInC))
	var revRegDefInC *C.char = C.findy_null_string
	if revRegDef != findy.NullString {
		revRegDefInC = C.CString(revRegDef)
		defer C.free(unsafe.Pointer(revRegDefInC))
	}
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyProverStoreCredential")
	C.findy_prover_store_credential(C.int(cmdHandle), C.int(wallet), credIDInC, credReqMetaInC, credentialsInC, credDefInC, revRegDefInC)
	return ch
}

func FindyProverSearchCredentialsForProofReq(wallet int, proofReqJSON, extraQueryJSON string) ctx.Channel {
	proofReqJSONInC := C.CString(proofReqJSON)
	defer C.free(unsafe.Pointer(proofReqJSONInC))
	var extraQueryJSONInC *C.char = C.findy_null_string
	if extraQueryJSON != findy.NullString {
		extraQueryJSONInC = C.CString(extraQueryJSON)
		defer C.free(unsafe.Pointer(extraQueryJSONInC))
	}
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyProverSearchCredentialsForProofReq")
	C.findy_prover_search_credentials_for_proof_req(C.int(cmdHandle), C.int(wallet), proofReqJSONInC, extraQueryJSONInC)
	return ch
}

func FindyProverFetchCredentialsForProofReq(searchHandle int, itemRef string, count int) ctx.Channel {
	itemRefInC := C.CString(itemRef)
	defer C.free(unsafe.Pointer(itemRefInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyProverFetchCredentialsForProofReq")
	C.findy_prover_fetch_credentials_for_proof_req(C.int(cmdHandle), C.int(searchHandle), itemRefInC, C.uint(count))
	return ch
}

func FindyProverCloseCredentialsSearchForProofReq(searchHandle int) ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyProverCloseCredentialsSearchForProofReq")
	C.findy_prover_close_credentials_search_for_proof_req(C.int(cmdHandle), C.int(searchHandle))
	return ch
}

func FindyVerifierVerifyProof(proofReqJSON, proofJSON, schemasJSON, credDefsJSON, revRegDefsJSON, revRegsJSON string) ctx.Channel {
	proofReqJSONInC := C.CString(proofReqJSON)
	defer C.free(unsafe.Pointer(proofReqJSONInC))
	proofJSONInC := C.CString(proofJSON)
	defer C.free(unsafe.Pointer(proofJSONInC))
	schemasJSONInC := C.CString(schemasJSON)
	defer C.free(unsafe.Pointer(schemasJSONInC))
	credDefsJSONInC := C.CString(credDefsJSON)
	defer C.free(unsafe.Pointer(credDefsJSONInC))
	revRegDefsJSONInC := C.CString(revRegDefsJSON)
	defer C.free(unsafe.Pointer(revRegDefsJSONInC))
	revRegsJSONInC := C.CString(revRegsJSON)
	defer C.free(unsafe.Pointer(revRegsJSONInC))
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyVerifierVerifyProof")
	C.findy_verifier_verify_proof(C.int(cmdHandle), proofReqJSONInC, proofJSONInC,
		schemasJSONInC, credDefsJSONInC, revRegDefsJSONInC, revRegsJSONInC)
	return ch
}
