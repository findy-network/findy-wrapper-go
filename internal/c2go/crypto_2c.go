package c2go

//#include <stdio.h>
//#include <stdlib.h>
//#include "findy_glue.h"
import "C"
import (
	"unsafe"

	"github.com/findy-network/findy-wrapper-go"
	"github.com/findy-network/findy-wrapper-go/internal/ctx"
)

// findy_crypto_sign
func FindyCryptoSign(wallet int, signerVerKey string, msg []byte) ctx.Channel {
	signerVerKeyInC := C.CString(signerVerKey)
	defer C.free(unsafe.Pointer(signerVerKeyInC))
	bytePtr := C.CBytes(msg)
	defer C.free(bytePtr)

	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyCryptoSign")
	C.findy_crypto_sign(C.int(cmdHandle), C.int(wallet), signerVerKeyInC, (*C.uchar)(bytePtr), C.uint(len(msg)))
	return ch
}

// findy_crypto_verify
func FindyCryptoVerify(signerVerKey string, rawMsg, rawSignature []byte) ctx.Channel {
	signerVerKeyInC := C.CString(signerVerKey)
	defer C.free(unsafe.Pointer(signerVerKeyInC))
	mBytes, sBytes := C.CBytes(rawMsg), C.CBytes(rawSignature)
	defer C.free(mBytes)
	defer C.free(sBytes)
	cmdHandle, ch := ctx.CmdContext.NamedPush("FindyCryptoVerify")
	C.findy_crypto_verify(C.int(cmdHandle), signerVerKeyInC, (*C.uchar)(mBytes),
		C.uint(len(rawMsg)), (*C.uchar)(sBytes), C.uint(len(rawSignature)))
	return ch
}

func CryptoAnonCrypt(recipientKey string, msgBytes []byte) ctx.Channel {
	recipientKeyInC := C.CString(recipientKey)
	bytePtr := C.CBytes(msgBytes)
	defer C.free(unsafe.Pointer(recipientKeyInC))
	defer C.free(bytePtr)
	cmdHandle, ch := ctx.CmdContext.NamedPush("AnonCrypt")
	C.findy_crypto_anon_crypt(C.int(cmdHandle), recipientKeyInC, (*C.uchar)(bytePtr), C.uint(len(msgBytes)))
	return ch
}

func CryptoAuthCrypt(wallet int, senderKey, recipientKey string, msgBytes []byte) ctx.Channel {
	senderKeyInC := C.CString(senderKey)
	recipientKeyInC := C.CString(recipientKey)
	bytePtr := C.CBytes(msgBytes)
	defer C.free(unsafe.Pointer(recipientKeyInC))
	defer C.free(unsafe.Pointer(senderKeyInC))
	defer C.free(bytePtr)
	cmdHandle, ch := ctx.CmdContext.NamedPush("AuthCrypt")
	C.findy_crypto_auth_crypt(C.int(cmdHandle), C.int(wallet), senderKeyInC,
		recipientKeyInC, (*C.uchar)(bytePtr), C.uint(len(msgBytes)))
	return ch
}

func CryptoAnonDecrypt(wallet int, recipientKey string, msgBytes []byte) ctx.Channel {
	recipientKeyInC := C.CString(recipientKey)
	bytePtr := C.CBytes(msgBytes)
	defer C.free(unsafe.Pointer(recipientKeyInC))
	defer C.free(bytePtr)
	cmdHandle, ch := ctx.CmdContext.NamedPush("AnonDecrypt")
	C.findy_crypto_anon_decrypt(C.int(cmdHandle), C.int(wallet), recipientKeyInC, (*C.uchar)(bytePtr), C.uint(len(msgBytes)))
	return ch
}

func CryptoAuthDecrypt(wallet int, recipientKey string, msgBytes []byte) ctx.Channel {
	recipientKeyInC := C.CString(recipientKey)
	bytePtr := C.CBytes(msgBytes)
	defer C.free(unsafe.Pointer(recipientKeyInC))
	defer C.free(bytePtr)
	cmdHandle, ch := ctx.CmdContext.NamedPush("AuthDecrypt")
	C.findy_crypto_auth_decrypt(C.int(cmdHandle), C.int(wallet), recipientKeyInC, (*C.uchar)(bytePtr), C.uint(len(msgBytes)))
	return ch
}

func FindyPackMessage(wallet int, msgBytes []byte, recipientKeysJSON, senderKey string) ctx.Channel {
	var senderKeyInC *C.char = C.findy_null_string
	if senderKey != findy.NullString {
		senderKeyInC = C.CString(senderKey)
		defer C.free(unsafe.Pointer(senderKeyInC))
	}
	recipientKeysJSONInC := C.CString(recipientKeysJSON)
	bytePtr := C.CBytes(msgBytes)
	defer C.free(unsafe.Pointer(recipientKeysJSONInC))
	defer C.free(bytePtr)

	cmdHandle, ch := ctx.CmdContext.NamedPush("PackMessage")
	C.findy_pack_message(C.int(cmdHandle), C.int(wallet), (*C.uchar)(bytePtr), C.uint(len(msgBytes)), recipientKeysJSONInC, senderKeyInC)
	return ch
}

func FindyUnpackMessage(wallet int, msgBytes []byte) ctx.Channel {
	bytePtr := C.CBytes(msgBytes)
	defer C.free(bytePtr)
	cmdHandle, ch := ctx.CmdContext.NamedPush("UnpackMessage")
	C.findy_unpack_message(C.int(cmdHandle), C.int(wallet), (*C.uchar)(bytePtr), C.uint(len(msgBytes)))
	return ch
}
