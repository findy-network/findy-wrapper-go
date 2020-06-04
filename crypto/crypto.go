/*
Package crypto is corresponding Go package for libindy's anoncreds namespace.
We suggest that you read indy SDK documentation for more information.
Unfortunately the documentation during the writing of this package was minimal
to nothing. We suggest you do the same as we did, read the rust code if more
detailed information is needed.
*/
package crypto

import (
	"github.com/optechlab/findy-go/dto"
	"github.com/optechlab/findy-go/internal/c2go"
	"github.com/optechlab/findy-go/internal/ctx"
)

// SignMsg signs a message with a verification key.
func SignMsg(wallet int, signerVerKey string, msg []byte) ctx.Channel {
	return c2go.FindyCryptoSign(wallet, signerVerKey, msg)
}

// VerifySignature verifies message's signature with a given verification key.
func VerifySignature(signerVerKey string, msg, sig []byte) ctx.Channel {
	return c2go.FindyCryptoVerify(signerVerKey, msg, sig)
}

// AnonCrypt encrypts the msg according the recipientKey
func AnonCrypt(recipientKey string, msg []byte) ctx.Channel {
	return c2go.CryptoAnonCrypt(recipientKey, msg)
}

// AnonDecrypt decrypts the msg according the recipientKey
func AnonDecrypt(wallet int, recipientKey string, msg []byte) ctx.Channel {
	return c2go.CryptoAnonDecrypt(wallet, recipientKey, msg)
}

// AuthCrypt encrypts and sings the msg.
func AuthCrypt(wallet int, senderKey, recipientKey string, msg []byte) ctx.Channel {
	return c2go.CryptoAuthCrypt(wallet, senderKey, recipientKey, msg)
}

// AuthDecrypt decrypts and verifies the msg.
func AuthDecrypt(wallet int, recipientKey string, msg []byte) ctx.Channel {
	return c2go.CryptoAuthDecrypt(wallet, recipientKey, msg)
}

// Pack encrypts and packs the original byte message to recipients and if
// sender key is set (use findy.NullString) signs it too. This is indy SDK Go
// wrapper.
func Pack(wallet int, senderKey string, msg []byte, recipientKeys ...string) ctx.Channel {
	return PackMessage(wallet, dto.DoJSONArray(recipientKeys), senderKey, msg)
}

// PackMessage encrypts and packs the original byte message to recipients and if
// sender key is set (use findy.NullString) signs it too. This is indy SDK Go
// wrapper.
func PackMessage(wallet int, recipientKeysJSON, senderKey string, msg []byte) ctx.Channel {
	return c2go.FindyPackMessage(wallet, msg, recipientKeysJSON, senderKey)
}

// UnpackMessage is a Go wrapper for same name indy SDK function to decrypt
// messages packed with PackMessage.
func UnpackMessage(wallet int, msg []byte) ctx.Channel {
	return c2go.FindyUnpackMessage(wallet, msg)
}

// Unpacked is a helper struct to wrap json type UnpackMessage function is
// returning.
type Unpacked struct {
	Message         string `json:"message"`
	RecipientVerkey string `json:"recipient_verkey"`
	SenderVerkey    string `json:"sender_verkey"`
}

// NewUnpacked creates a new instance of Unpacked helper wrapper of
// UnpackMessage function's result.
func NewUnpacked(msg []byte) *Unpacked {
	var unp Unpacked
	dto.FromJSON(msg, &unp)
	return &unp
}

// Bytes returns the actual decrypted message inside the Unpacked structure and
// inside of Message field.
func (msg *Unpacked) Bytes() []byte {
	return []byte(msg.Message)
}

// Unpacking is lazy fetch helper for UnpackMessage async method.
type Unpacking struct {
	Unpacked
	ch ctx.Channel
}

// NewUnpacking inits Unpacking struct correctly i.e., it calls UnpackMessage
// which is async method. This makes Unpacking a lazy fetch helper. The result
// will be fetched with Bytes() method.
func NewUnpacking(w int, msg []byte) *Unpacking {
	return &Unpacking{ch: UnpackMessage(w, msg)}
}

// Bytes returns previously started unpacking result bytes. Note! it's mandatory
// first call NewUnpacking() to start unpacking correctly.
func (up *Unpacking) Bytes() (bytes []byte, err error) {
	if up.ch == nil {
		panic("not initialized correctly")
	}

	r := <-up.ch
	if r.Err() != nil {
		return nil, r.Err()
	}

	up.Unpacked = *NewUnpacked(r.Bytes())
	return up.Unpacked.Bytes(), nil
}
