/*
Package findy is a Go wrapper for libindy, indy SDK. It follows libindy's naming
and sub package structure. The callback mechanism of the indy is changed to Go's
channel. All of the wrapper functions return the Channel. The input parameters
of the wrapper functions follow the libindy as well. They use same JSON objects.
However, you doesn't need to give JSON strings as an input arguments but similar
Go structures.

# About The Documentation

The Go documentation is minimal. We considered to paste libindy's documentation
from Rust source files but thought that was not necessary. We suggest you to
read indy SKD's documentation in cases the documentation in this wrapper is not
enough.

# Return Types

The return type is Channel, which transfers dto.Result Go structures. These
structures work similarly to C unions, which means that we can use one type for
all of the wrapper functions. To access the actual data, you have to know the actual
type. pool.OpenLedger() returns Handle type.

	r = <-pool.OpenLedger("FINDY_MEM_LEDGER")
	assert.NoError(r.Err())
	h2 := r.Handle()
	assert.Equal(h2, -1)

did.CreateAndStore() returns two strings: did and verkey. Please note the use of
did.Did struct instead of the JSON string.

	r := <-did.CreateAndStore(w, did.Did{Seed: ""})
	assert.NoError(r.Err())
	did := r.Str1()
	vk := r.Str2()

When a null string is needed for an argument, the predefined type must be used.

	r = <-annoncreds.IssuerCreateAndStoreCredentialDef(w1, w1DID, scJSON,
		"MY_FIRMS_CRED_DEF", findy.NullString, findy.NullString)
*/
package findy

import (
	"github.com/findy-network/findy-wrapper-go/internal/ctx"
)

// NullString is constant to pass null strings to the packet.
const NullString = "\xff"

// NullHandle is constant to pass null handles to the packet.
const NullHandle = -1

// Channel is channel type for findy API. Instead of callbacks findy API returns
// channels for async functions.
type Channel = ctx.Channel

// SetTrace enables or disables the trace output. It's disabled as default.
// Note, this is obsolete, use logging V style parameter, level 10
func SetTrace(_ bool) bool {
	return false
}
